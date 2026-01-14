package theme

import (
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color scheme for the TUI
type Theme struct {
	Name   string
	IsDark bool

	// Base colors
	Background lipgloss.Color
	Foreground lipgloss.Color

	// Accent colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color

	// Status colors
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color
	Muted   lipgloss.Color

	// Component-specific
	Selected    lipgloss.Color
	Cursor      lipgloss.Color
	Border      lipgloss.Color
	HeaderBg    lipgloss.Color
	StatusBarBg lipgloss.Color
}

// DarkTheme is the default dark color scheme (Tokyo Night inspired)
var DarkTheme = Theme{
	Name:        "dark",
	IsDark:      true,
	Background:  lipgloss.Color("#1a1b26"),
	Foreground:  lipgloss.Color("#c0caf5"),
	Primary:     lipgloss.Color("#7aa2f7"),
	Secondary:   lipgloss.Color("#bb9af7"),
	Accent:      lipgloss.Color("#7dcfff"),
	Success:     lipgloss.Color("#9ece6a"),
	Warning:     lipgloss.Color("#e0af68"),
	Error:       lipgloss.Color("#f7768e"),
	Muted:       lipgloss.Color("#565f89"),
	Selected:    lipgloss.Color("#33467c"),
	Cursor:      lipgloss.Color("#7aa2f7"),
	Border:      lipgloss.Color("#3b4261"),
	HeaderBg:    lipgloss.Color("#24283b"),
	StatusBarBg: lipgloss.Color("#1f2335"),
}

// LightTheme is the light color scheme
var LightTheme = Theme{
	Name:        "light",
	IsDark:      false,
	Background:  lipgloss.Color("#f5f5f5"),
	Foreground:  lipgloss.Color("#343b58"),
	Primary:     lipgloss.Color("#34548a"),
	Secondary:   lipgloss.Color("#5a4a78"),
	Accent:      lipgloss.Color("#166775"),
	Success:     lipgloss.Color("#485e30"),
	Warning:     lipgloss.Color("#8c6c3e"),
	Error:       lipgloss.Color("#8c4351"),
	Muted:       lipgloss.Color("#9699a3"),
	Selected:    lipgloss.Color("#d4d6e4"),
	Cursor:      lipgloss.Color("#34548a"),
	Border:      lipgloss.Color("#c0c0c0"),
	HeaderBg:    lipgloss.Color("#e8e8e8"),
	StatusBarBg: lipgloss.Color("#e0e0e0"),
}

// Current holds the active theme
var Current = DarkTheme

// DetectSystemTheme attempts to detect the system's color scheme
func DetectSystemTheme() Theme {
	// Check COLORFGBG environment variable (format: fg;bg)
	if fgbg := os.Getenv("COLORFGBG"); fgbg != "" {
		parts := strings.Split(fgbg, ";")
		if len(parts) >= 2 {
			bg, err := strconv.Atoi(parts[len(parts)-1])
			if err == nil && bg >= 0 && bg < 7 {
				return DarkTheme
			}
			return LightTheme
		}
	}

	// Check for common dark terminal indicators
	if term := os.Getenv("TERM"); strings.Contains(term, "256color") {
		// Most 256-color terminals default to dark
		return DarkTheme
	}

	// Default to dark theme
	return DarkTheme
}

// Toggle switches between dark and light themes
func Toggle() Theme {
	if Current.IsDark {
		Current = LightTheme
	} else {
		Current = DarkTheme
	}
	return Current
}

// Set sets the current theme
func Set(t Theme) {
	Current = t
}

// Styles returns lipgloss styles for the current theme
type Styles struct {
	// Text styles
	Title   lipgloss.Style
	Normal  lipgloss.Style
	Muted   lipgloss.Style
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style

	// Component styles
	Header      lipgloss.Style
	StatusBar   lipgloss.Style
	Selected    lipgloss.Style
	Cursor      lipgloss.Style
	Border      lipgloss.Style
	Panel       lipgloss.Style
	Dialog      lipgloss.Style
	DialogTitle lipgloss.Style
}

// GetStyles returns styles based on the current theme
func GetStyles() Styles {
	t := Current

	return Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary),

		Normal: lipgloss.NewStyle().
			Foreground(t.Foreground),

		Muted: lipgloss.NewStyle().
			Foreground(t.Muted),

		Success: lipgloss.NewStyle().
			Foreground(t.Success),

		Warning: lipgloss.NewStyle().
			Foreground(t.Warning),

		Error: lipgloss.NewStyle().
			Foreground(t.Error),

		Header: lipgloss.NewStyle().
			Background(t.HeaderBg).
			Foreground(t.Foreground).
			Padding(0, 1),

		StatusBar: lipgloss.NewStyle().
			Background(t.StatusBarBg).
			Foreground(t.Muted).
			Padding(0, 1),

		Selected: lipgloss.NewStyle().
			Background(t.Selected).
			Foreground(t.Foreground).
			Bold(true),

		Cursor: lipgloss.NewStyle().
			Foreground(t.Cursor).
			Bold(true),

		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border),

		Panel: lipgloss.NewStyle().
			Padding(0, 1),

		Dialog: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Padding(1, 2),

		DialogTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary).
			Padding(0, 1),
	}
}
