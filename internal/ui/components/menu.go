package components

import (
	"fmt"
	"strings"

	"ccm/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

// MenuModel is a simple action menu
type MenuModel struct {
	label    string
	items    []string
	cursor   int
	selected int
	quitting bool
	canceled bool
}

// NewMenu creates a new action menu
func NewMenu(items []string, label string) MenuModel {
	return MenuModel{
		label:    label,
		items:    items,
		cursor:   0,
		selected: -1,
	}
}

// Init implements tea.Model
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c", "esc", "q":
			m.canceled = true
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model
func (m MenuModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Label
	b.WriteString(styles.PromptStyle.Render(m.label))
	b.WriteString("\n")

	// Items
	for i, item := range m.items {
		cursor := "  "
		style := styles.NormalItemStyle
		if i == m.cursor {
			cursor = styles.CursorStyle.Render("▸ ")
			style = styles.SelectedItemStyle
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(item)))
	}

	// Help
	b.WriteString(styles.MutedStyle.Render("\n↑/↓: move • enter: select • q: quit"))

	return b.String()
}

// Selected returns the selected index
func (m MenuModel) Selected() int {
	return m.selected
}

// Canceled returns whether the user canceled
func (m MenuModel) Canceled() bool {
	return m.canceled
}
