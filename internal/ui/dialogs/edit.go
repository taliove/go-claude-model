package dialogs

import (
	"fmt"
	"strings"

	"ccm/internal/provider"
	"ccm/internal/ui/messages"
	"ccm/internal/ui/theme"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EditDialogModel is the edit provider dialog
type EditDialogModel struct {
	provider     provider.Provider
	fields       []textinput.Model
	focusedField int
	width        int
	height       int
	submitted    bool
	canceled     bool
}

const (
	fieldAPIKey = iota
	fieldBaseURL
	fieldModel
	numFields
)

// NewEditDialog creates a new edit dialog
func NewEditDialog(p provider.Provider) EditDialogModel {
	fields := make([]textinput.Model, numFields)

	// API Key field
	fields[fieldAPIKey] = textinput.New()
	fields[fieldAPIKey].Placeholder = "sk-..."
	fields[fieldAPIKey].EchoMode = textinput.EchoPassword
	fields[fieldAPIKey].EchoCharacter = '*'
	fields[fieldAPIKey].CharLimit = 256
	fields[fieldAPIKey].Width = 40
	if p.APIKey != "" {
		fields[fieldAPIKey].SetValue(p.APIKey)
	}

	// Base URL field
	fields[fieldBaseURL] = textinput.New()
	fields[fieldBaseURL].Placeholder = "https://api.example.com"
	fields[fieldBaseURL].CharLimit = 256
	fields[fieldBaseURL].Width = 40
	fields[fieldBaseURL].SetValue(p.BaseURL)

	// Model field
	fields[fieldModel] = textinput.New()
	fields[fieldModel].Placeholder = "model-name"
	fields[fieldModel].CharLimit = 128
	fields[fieldModel].Width = 40
	fields[fieldModel].SetValue(p.Model)

	// Focus first field
	fields[fieldAPIKey].Focus()

	return EditDialogModel{
		provider:     p,
		fields:       fields,
		focusedField: 0,
		width:        60,
		height:       14,
	}
}

// Title returns the dialog title
func (m EditDialogModel) Title() string {
	return fmt.Sprintf("Edit Provider: %s", m.provider.Name)
}

// Width returns the dialog width
func (m EditDialogModel) Width() int {
	return m.width
}

// Height returns the dialog height
func (m EditDialogModel) Height() int {
	return m.height
}

// Submitted returns whether the dialog was submitted
func (m EditDialogModel) Submitted() bool {
	return m.submitted
}

// Canceled returns whether the dialog was canceled
func (m EditDialogModel) Canceled() bool {
	return m.canceled
}

// GetProvider returns the updated provider
func (m EditDialogModel) GetProvider() provider.Provider {
	p := m.provider
	p.APIKey = m.fields[fieldAPIKey].Value()
	p.BaseURL = m.fields[fieldBaseURL].Value()
	p.Model = m.fields[fieldModel].Value()
	return p
}

// Init implements tea.Model
func (m EditDialogModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m EditDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.fields[m.focusedField].Blur()
			m.focusedField = (m.focusedField + 1) % numFields
			m.fields[m.focusedField].Focus()
			return m, textinput.Blink
		case "shift+tab", "up":
			m.fields[m.focusedField].Blur()
			m.focusedField = (m.focusedField - 1 + numFields) % numFields
			m.fields[m.focusedField].Focus()
			return m, textinput.Blink
		case "enter":
			// Validate
			if m.fields[fieldAPIKey].Value() == "" {
				return m, nil
			}
			m.submitted = true
			return m, func() tea.Msg {
				return messages.CloseDialogMsg{Result: m.GetProvider()}
			}
		case "esc":
			m.canceled = true
			return m, func() tea.Msg {
				return messages.CloseDialogMsg{Result: nil}
			}
		}
	}

	// Update focused field
	var cmd tea.Cmd
	m.fields[m.focusedField], cmd = m.fields[m.focusedField].Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m EditDialogModel) View() string {
	styles := theme.GetStyles()
	t := theme.Current

	var b strings.Builder

	// Title
	title := styles.DialogTitle.Render(m.Title())
	b.WriteString(title)
	b.WriteString("\n\n")

	// Provider info
	providerInfo := lipgloss.NewStyle().
		Foreground(t.Muted).
		Render(fmt.Sprintf("(%s)", m.provider.DisplayName))
	b.WriteString(providerInfo)
	b.WriteString("\n\n")

	// Fields
	fieldNames := []string{"API Key:", "Base URL:", "Model:"}
	labelStyle := lipgloss.NewStyle().
		Width(10).
		Foreground(t.Foreground)

	for i, name := range fieldNames {
		isFocused := i == m.focusedField

		label := labelStyle.Render(name)
		field := m.fields[i].View()

		if isFocused {
			label = lipgloss.NewStyle().
				Width(10).
				Foreground(t.Primary).
				Bold(true).
				Render(name)
		}

		b.WriteString(label)
		b.WriteString(" ")
		b.WriteString(field)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Buttons hint
	buttonHint := styles.Muted.Render("Tab: next field  Enter: save  Esc: cancel")
	b.WriteString(buttonHint)

	// Wrap in dialog box
	content := b.String()
	return styles.Dialog.Width(m.width).Render(content)
}
