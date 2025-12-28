package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
)

// AppState represents the current state of the TUI application
type AppState int

const (
	StateWelcome AppState = iota
	StateEditor
	StateOutput
	StateLoading
)

// Model is the main TUI application model
type Model struct {
	state   AppState
	mode    string
	width   int
	height  int
	err     error
	welcome *Welcome
	editor  *EditorModel
}

// ErrorMsg is sent when an error occurs
type ErrorMsg struct {
	err error
}

// EnterEditorMsg is sent to transition to the editor state
type EnterEditorMsg struct{}

// MessageModel represents a chat message in the conversation
type MessageModel struct {
	Role      string   `json:"role"`
	Content   string   `json:"content"`
	Metadata  string   `json:"metadata,omitempty"`
	Questions []string `json:"questions,omitempty"`
}

// EditorModel is the code editor interface
type EditorModel struct {
	input               textinput.Model
	viewport            viewport.Model
	width               int
	height              int
	conversationHistory []MessageModel
	isFetching          bool
	spinner             spinner.Model
	glamourRenderer     *glamour.TermRenderer
	ready               bool
	shouldScrollBottom  bool
}

// AgentResponseMsg is sent when the agent completes a request
type AgentResponseMsg struct {
	userQuery string
	code      string
	metadata  string
	questions []string
}

// AgentErrorMsg is sent when the agent encounters an error
type AgentErrorMsg struct {
	userQuery string
	err       error
}

// Welcome is the welcome screen model
type Welcome struct {
	mode     string
	input    string
	commands []Command
}

// Command represents an available TUI command
type Command struct {
	Name        string
	Description string
	Available   bool
}

// ServerStartedMsg is sent when the server starts
type ServerStartedMsg struct{}

// IngesterCompleteMsg is sent when the ingester completes
type IngesterCompleteMsg struct{}
