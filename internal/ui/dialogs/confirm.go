package dialogs

import (
	"fmt"
	"strings"

	"ccm/internal/ui/messages"
	"ccm/internal/ui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmDialogModel is the confirmation dialog
type ConfirmDialogModel struct {
	title    string
	message  string
	selected int // 0=yes, 1=no
	width    int
	height   int
	result   bool
	answered bool
}

// NewConfirmDialog creates a new confirmation dialog
func NewConfirmDialog(title, message string) ConfirmDialogModel {
	return ConfirmDialogModel{
		title:    title,
		message:  message,
		selected: 1, // Default to "No"
		width:    50,
		height:   8,
	}
}

// Title returns the dialog title
func (m ConfirmDialogModel) Title() string {
	return m.title
}

// Width returns the dialog width
func (m ConfirmDialogModel) Width() int {
	return m.width
}

// Height returns the dialog height
func (m ConfirmDialogModel) Height() int {
	return m.height
}

// Result returns the confirmation result
func (m ConfirmDialogModel) Result() bool {
	return m.result
}

// Answered returns whether the user answered
func (m ConfirmDialogModel) Answered() bool {
	return m.answered
}

// Init implements tea.Model
func (m ConfirmDialogModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ConfirmDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.selected = 0
		case "right", "l":
			m.selected = 1
		case "tab":
			m.selected = (m.selected + 1) % 2
		case "y", "Y":
			m.result = true
			m.answered = true
			return m, func() tea.Msg {
				return messages.CloseDialogMsg{Result: true}
			}
		case "n", "N", "esc":
			m.result = false
			m.answered = true
			return m, func() tea.Msg {
				return messages.CloseDialogMsg{Result: false}
			}
		case "enter":
			m.result = m.selected == 0
			m.answered = true
			return m, func() tea.Msg {
				return messages.CloseDialogMsg{Result: m.result}
			}
		}
	}
	return m, nil
}

// View implements tea.Model
func (m ConfirmDialogModel) View() string {
	styles := theme.GetStyles()
	t := theme.Current

	var b strings.Builder

	// Title
	title := styles.DialogTitle.Render(m.title)
	b.WriteString(title)
	b.WriteString("\n\n")

	// Message
	msgStyle := lipgloss.NewStyle().Foreground(t.Foreground)
	b.WriteString(msgStyle.Render(m.message))
	b.WriteString("\n\n")

	// Buttons
	baseStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Muted)
	yesStyle := baseStyle
	noStyle := baseStyle

	if m.selected == 0 {
		yesStyle = yesStyle.
			BorderForeground(t.Primary).
			Foreground(t.Primary).
			Bold(true)
	} else {
		noStyle = noStyle.
			BorderForeground(t.Primary).
			Foreground(t.Primary).
			Bold(true)
	}

	yesBtn := yesStyle.Render("Yes")
	noBtn := noStyle.Render("No")

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, yesBtn, "  ", noBtn)
	b.WriteString(buttons)
	b.WriteString("\n\n")

	// Hint
	hint := styles.Muted.Render("y: yes  n: no  Tab: switch  Enter: confirm")
	b.WriteString(hint)

	// Center buttons
	content := b.String()
	return styles.Dialog.Width(m.width).Render(content)
}

// RemoveConfirmDialog creates a dialog for removing a provider
func RemoveConfirmDialog(name, displayName string) ConfirmDialogModel {
	return NewConfirmDialog(
		"Remove Provider",
		fmt.Sprintf("Remove provider '%s' (%s)?\nThis will delete the configuration.", name, displayName),
	)
}
