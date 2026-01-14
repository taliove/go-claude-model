package components

import (
	"fmt"

	"ccm/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmModel is a Y/N confirmation dialog
type ConfirmModel struct {
	label     string
	confirmed bool
	answered  bool
	canceled  bool
}

// NewConfirm creates a new confirmation dialog
func NewConfirm(label string) ConfirmModel {
	return ConfirmModel{
		label: label,
	}
}

// Init implements tea.Model
func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
			m.answered = true
			return m, tea.Quit
		case "n", "N", "enter":
			m.confirmed = false
			m.answered = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.canceled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model
func (m ConfirmModel) View() string {
	if m.answered || m.canceled {
		return ""
	}
	prompt := styles.PromptStyle.Render(m.label)
	hint := styles.MutedStyle.Render(" [y/N] ")
	return fmt.Sprintf("%s%s", prompt, hint)
}

// Confirmed returns whether the user confirmed
func (m ConfirmModel) Confirmed() bool {
	return m.confirmed
}

// Canceled returns whether the user canceled
func (m ConfirmModel) Canceled() bool {
	return m.canceled
}
