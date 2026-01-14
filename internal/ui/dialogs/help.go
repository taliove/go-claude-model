package dialogs

import (
	"strings"

	"ccm/internal/ui/messages"
	"ccm/internal/ui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpDialogModel shows keyboard shortcuts
type HelpDialogModel struct {
	width  int
	height int
}

// NewHelpDialog creates a new help dialog
func NewHelpDialog() HelpDialogModel {
	return HelpDialogModel{
		width:  60,
		height: 20,
	}
}

// Title returns the dialog title
func (m HelpDialogModel) Title() string {
	return "Keyboard Shortcuts"
}

// Width returns the dialog width
func (m HelpDialogModel) Width() int {
	return m.width
}

// Height returns the dialog height
func (m HelpDialogModel) Height() int {
	return m.height
}

// Init implements tea.Model
func (m HelpDialogModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m HelpDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Any key closes the help dialog
		return m, func() tea.Msg {
			return messages.CloseDialogMsg{}
		}
	}
	return m, nil
}

// View implements tea.Model
func (m HelpDialogModel) View() string {
	styles := theme.GetStyles()
	t := theme.Current

	var b strings.Builder

	// Title
	title := styles.DialogTitle.Render(m.Title())
	b.WriteString(title)
	b.WriteString("\n\n")

	keyStyle := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Width(12)
	descStyle := lipgloss.NewStyle().
		Foreground(t.Foreground)
	sectionStyle := lipgloss.NewStyle().
		Foreground(t.Warning).
		Bold(true)

	// Navigation section
	b.WriteString(sectionStyle.Render("Navigation"))
	b.WriteString("\n")
	shortcuts := []struct {
		key  string
		desc string
	}{
		{"j/k, ↑/↓", "Move cursor up/down"},
		{"g/G", "Jump to first/last"},
		{"/", "Start search"},
		{"Esc", "Cancel search / Close dialog"},
	}
	for _, s := range shortcuts {
		b.WriteString(keyStyle.Render(s.key))
		b.WriteString(descStyle.Render(s.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Actions section
	b.WriteString(sectionStyle.Render("Actions"))
	b.WriteString("\n")
	actions := []struct {
		key  string
		desc string
	}{
		{"Enter", "Run Claude with selected provider"},
		{"e", "Edit provider configuration"},
		{"d", "Set as default provider"},
		{"t", "Test provider connection"},
		{"r", "Remove provider"},
		{"a", "Add/configure provider"},
	}
	for _, s := range actions {
		b.WriteString(keyStyle.Render(s.key))
		b.WriteString(descStyle.Render(s.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// General section
	b.WriteString(sectionStyle.Render("General"))
	b.WriteString("\n")
	general := []struct {
		key  string
		desc string
	}{
		{"?", "Toggle this help"},
		{"Ctrl+T", "Toggle dark/light theme"},
		{"q", "Quit"},
	}
	for _, s := range general {
		b.WriteString(keyStyle.Render(s.key))
		b.WriteString(descStyle.Render(s.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.Muted.Render("Press any key to close"))

	content := b.String()
	return styles.Dialog.Width(m.width).Render(content)
}
