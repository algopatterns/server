package ssh

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/algorave/server/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
)

const (
	defaultPort    = 2222
	defaultTimeout = 30 * time.Minute
)

type Config struct {
	Port           int
	HostKeyPath    string
	MaxConnections int
	IdleTimeout    time.Duration
	ProductionMode bool
}

// returns default SSH configuration
func DefaultConfig() *Config {
	return &Config{
		Port:           defaultPort,
		HostKeyPath:    ".ssh/algorave_host_key",
		MaxConnections: 50,
		IdleTimeout:    defaultTimeout,
		ProductionMode: true,
	}
}

type Server struct {
	config *Config
	server *ssh.Server
}

// returns a new SSH server
func NewServer(cfg *Config) (*Server, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	s := &Server{config: cfg}

	server, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf(":%d", cfg.Port)),
		wish.WithHostKeyPath(cfg.HostKeyPath),
		wish.WithMiddleware(
			bm.Middleware(s.tuiHandler),
			lm.Middleware(),
		),
		wish.WithIdleTimeout(cfg.IdleTimeout),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create SSH server: %w", err)
	}

	s.server = server
	return s, nil
}

// returns a TUI instance for each SSH connection
func (s *Server) tuiHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	mode := "production"
	if !s.config.ProductionMode {
		mode = "development"
	}

	app := tui.NewApp(mode)
	opts := []tea.ProgramOption{tea.WithAltScreen()}

	return app, opts
}

// starts the SSH server
func (s *Server) Start() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("🔐 SSH server starting on port %d", s.config.Port)
	log.Printf("   mode: %s (guest access)", modeString(s.config.ProductionMode))
	log.Printf("   connect: ssh localhost -p %d", s.config.Port)

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatalf("SSH server failed: %v", err)
		}
	}()

	<-done

	log.Println("🛑 shutting down SSH server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("SSH shutdown failed: %w", err)
	}

	log.Println("✓ SSH server exited gracefully")
	return nil
}

func modeString(production bool) string {
	if production {
		return "production"
	}

	return "development"
}
