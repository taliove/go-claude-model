package components

import (
	"ccm/internal/ui/theme"

	"github.com/charmbracelet/lipgloss"
)

// StatusBarModel represents the bottom status bar
type StatusBarModel struct {
	width       int
	message     string
	isError     bool
	defaultName string
	themeIcon   string
}

// NewStatusBar creates a new status bar
func NewStatusBar() StatusBarModel {
	return StatusBarModel{
		themeIcon: "◐", // half-filled circle for theme indicator
	}
}

// SetWidth updates the status bar width
func (m *StatusBarModel) SetWidth(width int) {
	m.width = width
}

// SetMessage sets a temporary message
func (m *StatusBarModel) SetMessage(msg string, isError bool) {
	m.message = msg
	m.isError = isError
}

// ClearMessage clears the temporary message
func (m *StatusBarModel) ClearMessage() {
	m.message = ""
	m.isError = false
}

// SetDefaultProvider sets the default provider name
func (m *StatusBarModel) SetDefaultProvider(name string) {
	m.defaultName = name
}

// SetThemeIcon updates the theme indicator
func (m *StatusBarModel) SetThemeIcon(isDark bool) {
	if isDark {
		m.themeIcon = "◐" // dark mode
	} else {
		m.themeIcon = "◑" // light mode
	}
}

// View renders the status bar
func (m StatusBarModel) View() string {
	styles := theme.GetStyles()
	t := theme.Current

	// Left side: shortcuts
	shortcuts := []struct {
		key  string
		desc string
	}{
		{"Enter", "run"},
		{"e", "edit"},
		{"d", "default"},
		{"t", "test"},
		{"r", "remove"},
		{"/", "search"},
		{"?", "help"},
		{"q", "quit"},
	}

	keyStyle := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(t.Muted)

	var shortcutText string
	for i, s := range shortcuts {
		if i > 0 {
			shortcutText += "  "
		}
		shortcutText += keyStyle.Render(s.key) + descStyle.Render(":"+s.desc)
	}

	// Right side: default provider and theme
	rightContent := ""
	if m.defaultName != "" {
		defaultStyle := lipgloss.NewStyle().
			Foreground(t.Warning)
		rightContent += defaultStyle.Render("★ " + m.defaultName)
		rightContent += "  "
	}
	rightContent += lipgloss.NewStyle().Foreground(t.Muted).Render(m.themeIcon)

	// Calculate spacing
	leftWidth := lipgloss.Width(shortcutText)
	rightWidth := lipgloss.Width(rightContent)
	spacing := m.width - leftWidth - rightWidth - 4 // padding

	if spacing < 1 {
		spacing = 1
	}

	spacer := lipgloss.NewStyle().Width(spacing).Render("")

	// If there's a message, show it instead of shortcuts
	var content string
	if m.message != "" {
		msgStyle := styles.Normal
		if m.isError {
			msgStyle = styles.Error
		}
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			msgStyle.Render(m.message),
			spacer,
			rightContent,
		)
	} else {
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			shortcutText,
			spacer,
			rightContent,
		)
	}

	return styles.StatusBar.Width(m.width).Render(content)
}
