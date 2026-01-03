package strudel

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type ValidationResult struct {
	Valid  bool   `json:"valid"`
	Error  string `json:"error,omitempty"`
	Line   *int   `json:"line,omitempty"`
	Column *int   `json:"column,omitempty"`
}

type validatorRequest struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Code string `json:"code,omitempty"`
}

type validatorResponse struct {
	ID     string `json:"id"`
	Ready  bool   `json:"ready,omitempty"`
	Valid  bool   `json:"valid,omitempty"`
	Error  string `json:"error,omitempty"`
	Line   *int   `json:"line,omitempty"`
	Column *int   `json:"column,omitempty"`
	Pong   bool   `json:"pong,omitempty"`
}

type Validator struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    *bufio.Scanner
	mu        sync.Mutex
	requestID atomic.Uint64
	ready     bool
	scriptDir string
}

// creates a new Strudel code validator.
// scriptDir is the path to the validate-strudel script directory.
func NewValidator(scriptDir string) (*Validator, error) {
	v := &Validator{
		scriptDir: scriptDir,
	}

	if err := v.start(); err != nil {
		return nil, err
	}

	return v, nil
}

// creates a validator using the server root directory.
func NewValidatorFromRoot(serverRoot string) (*Validator, error) {
	scriptDir := filepath.Join(serverRoot, "scripts", "validate-strudel")
	return NewValidator(scriptDir)
}

func (v *Validator) start() error {
	// try to find a compiled binary first, fall back to node
	cmd, err := v.findValidatorCommand()
	if err != nil {
		return err
	}

	v.cmd = cmd
	v.cmd.Dir = v.scriptDir

	stdin, err := v.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	v.stdin = stdin

	stdout, err := v.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	v.stdout = bufio.NewScanner(stdout)

	if err := v.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start validator process: %w", err)
	}

	// wait for ready signal
	resp, err := v.readResponse()
	if err != nil {
		return fmt.Errorf("validator process closed unexpectedly: %w", err)
	}

	if resp.Error != "" {
		return fmt.Errorf("validator error: %s", resp.Error)
	}

	if !resp.Ready {
		return fmt.Errorf("expected ready signal, got: %+v", resp)
	}

	v.ready = true
	return nil
}

// checks if Strudel code is syntactically valid.
func (v *Validator) Validate(ctx context.Context, code string) (*ValidationResult, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.ready {
		return nil, fmt.Errorf("validator not ready")
	}

	id := fmt.Sprintf("%d", v.requestID.Add(1))

	req := validatorRequest{
		Type: "validate",
		ID:   id,
		Code: code,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if _, err := v.stdin.Write(append(reqBytes, '\n')); err != nil {
		return nil, fmt.Errorf("failed to write to validator: %w", err)
	}

	// read response with timeout
	type result struct {
		resp *validatorResponse
		err  error
	}
	done := make(chan result, 1)

	go func() {
		resp, err := v.readResponse()
		done <- result{resp, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("validation timeout")
	case r := <-done:
		if r.err != nil {
			return nil, r.err
		}
		if r.resp.ID != id {
			return nil, fmt.Errorf("response ID mismatch: expected %s, got %s", id, r.resp.ID)
		}
		return &ValidationResult{
			Valid:  r.resp.Valid,
			Error:  r.resp.Error,
			Line:   r.resp.Line,
			Column: r.resp.Column,
		}, nil
	}
}

// shuts down validator process.
func (v *Validator) Close() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.stdin != nil {
		v.stdin.Close() //nolint:errcheck,gosec // best-effort cleanup
	}

	if v.cmd != nil && v.cmd.Process != nil {
		return v.cmd.Process.Kill()
	}

	return nil
}

// returns whether validator is ready to accept requests.
func (v *Validator) IsReady() bool {
	return v.ready
}

// reads the next JSON response from stdout, skipping non-JSON lines.
func (v *Validator) readResponse() (*validatorResponse, error) {
	for {
		if !v.stdout.Scan() {
			return nil, fmt.Errorf("validator closed or scan error")
		}

		line := v.stdout.Bytes()

		// skip non-JSON lines (e.g., strudel's emoji banner)
		if len(line) == 0 || line[0] != '{' {
			continue
		}

		var resp validatorResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			continue // skip malformed JSON
		}

		return &resp, nil
	}
}

// returns the command to run the validator.
func (v *Validator) findValidatorCommand() (*exec.Cmd, error) {
	var binaryName string

	switch runtime.GOOS {
	case "linux":
		binaryName = "validator-linuxstatic-x64"
	case "darwin":
		if runtime.GOARCH == "arm64" {
			binaryName = "validator-macos-arm64"
		} else {
			binaryName = "validator-macos-x64"
		}
	}

	// check for compiled binary in dist/
	if binaryName != "" {
		binaryPath := filepath.Join(v.scriptDir, "dist", binaryName)

		if _, err := os.Stat(binaryPath); err == nil {
			return exec.Command(binaryPath), nil //nolint:gosec // path is constructed from known components
		}

		// also check directly in scriptDir (for Docker deployments)
		binaryPath = filepath.Join(v.scriptDir, binaryName)

		if _, err := os.Stat(binaryPath); err == nil {
			return exec.Command(binaryPath), nil //nolint:gosec // path is constructed from known components
		}
	}

	// fall back to node
	scriptPath := filepath.Join(v.scriptDir, "validator.js")

	if _, err := os.Stat(scriptPath); err != nil {
		return nil, fmt.Errorf("validator not found: no binary for %s/%s and no validator.js", runtime.GOOS, runtime.GOARCH)
	}

	return exec.Command("node", scriptPath), nil //nolint:gosec // path is validated above
}
