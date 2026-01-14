package components

import (
	"fmt"

	"ccm/internal/ui/theme"
	"ccm/internal/version"

	"github.com/charmbracelet/lipgloss"
)

// HeaderModel represents the header bar
type HeaderModel struct {
	width int
}

// NewHeader creates a new header model
func NewHeader() HeaderModel {
	return HeaderModel{}
}

// SetWidth updates the header width
func (m *HeaderModel) SetWidth(width int) {
	m.width = width
}

// View renders the header
func (m HeaderModel) View() string {
	styles := theme.GetStyles()
	t := theme.Current

	title := styles.Title.Render("CCM - Claude Code Manager")
	ver := lipgloss.NewStyle().
		Foreground(t.Muted).
		Render(fmt.Sprintf("v%s", version.Version))

	// Calculate spacing
	titleWidth := lipgloss.Width(title)
	verWidth := lipgloss.Width(ver)
	spacing := m.width - titleWidth - verWidth - 2 // 2 for padding

	if spacing < 1 {
		spacing = 1
	}

	spacer := lipgloss.NewStyle().Width(spacing).Render("")

	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		spacer,
		ver,
	)

	return styles.Header.Width(m.width).Render(headerContent)
}
