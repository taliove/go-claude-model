package components

import (
	"fmt"

	"ccm/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TextInputModel is a text input with optional masking
type TextInputModel struct {
	textInput textinput.Model
	label     string
	err       error
	submitted bool
	canceled  bool
}

// NewTextInput creates a new text input
func NewTextInput(label string, masked bool) TextInputModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	if masked {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '*'
	}

	return TextInputModel{
		textInput: ti,
		label:     label,
	}
}

// NewAPIKeyInput creates a masked input for API keys
func NewAPIKeyInput(label string) TextInputModel {
	m := NewTextInput(label, true)
	m.textInput.Placeholder = "sk-..."
	return m
}

// Init implements tea.Model
func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.textInput.Value() == "" {
				m.err = fmt.Errorf("input cannot be empty")
				return m, nil
			}
			m.submitted = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.canceled = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m TextInputModel) View() string {
	if m.submitted || m.canceled {
		return ""
	}

	prompt := styles.PromptStyle.Render(m.label + ": ")
	input := m.textInput.View()

	result := fmt.Sprintf("%s%s", prompt, input)

	if m.err != nil {
		result += "\n" + styles.ErrorStyle.Render(m.err.Error())
	}

	result += "\n" + styles.MutedStyle.Render("enter: confirm • esc: cancel")

	return result
}

// Value returns the input value
func (m TextInputModel) Value() string {
	return m.textInput.Value()
}

// Submitted returns whether the input was submitted
func (m TextInputModel) Submitted() bool {
	return m.submitted
}

// Canceled returns whether the user canceled
func (m TextInputModel) Canceled() bool {
	return m.canceled
}
