package styles

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	Primary   = lipgloss.Color("#00BFFF")
	Success   = lipgloss.Color("#00FF00")
	Warning   = lipgloss.Color("#FFD700")
	Error     = lipgloss.Color("#FF0000")
	Muted     = lipgloss.Color("#808080")
	Highlight = lipgloss.Color("#FF00FF")
)

// Text styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)
)

// List item styles
var (
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	CursorStyle = lipgloss.NewStyle().
			Foreground(Primary)
)

// Status indicator styles
var (
	ConfiguredStyle = lipgloss.NewStyle().
			Foreground(Success)

	UnconfiguredStyle = lipgloss.NewStyle().
				Foreground(Error)

	DefaultMarkerStyle = lipgloss.NewStyle().
				Foreground(Warning)
)

// Input styles
var (
	PromptStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))
)
