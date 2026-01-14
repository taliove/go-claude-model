package app

import (
	"strings"

	"ccm/internal/ui/theme"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model
func (m AppModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.quitting {
		return ""
	}

	var sections []string

	// Header
	sections = append(sections, m.header.View())

	// Main content (provider list)
	listView := m.providerList.View()
	sections = append(sections, listView)

	// Detail panel (if provider selected)
	if m.providerList.Selected() != nil {
		sections = append(sections, m.detailPanel.View())
	}

	// Status bar
	sections = append(sections, m.statusBar.View())

	// Join sections
	content := strings.Join(sections, "\n")

	// If dialog is open, overlay it
	if m.activeDialog != nil {
		content = m.overlayDialog(content)
	}

	return content
}

// overlayDialog renders the dialog centered over the content
func (m AppModel) overlayDialog(content string) string {
	if m.activeDialog == nil {
		return content
	}

	dialogView := m.activeDialog.View()
	dialogWidth := m.activeDialog.Width()
	dialogHeight := m.activeDialog.Height()

	// Calculate center position
	x := (m.width - dialogWidth) / 2
	y := (m.height - dialogHeight) / 2

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Ensure we have enough lines
	for len(lines) < m.height {
		lines = append(lines, "")
	}

	// Split dialog into lines
	dialogLines := strings.Split(dialogView, "\n")

	// Overlay dialog
	for i, dialogLine := range dialogLines {
		lineIdx := y + i
		if lineIdx >= len(lines) {
			break
		}

		// Ensure the line is long enough
		line := lines[lineIdx]
		for lipgloss.Width(line) < m.width {
			line += " "
		}

		// Insert dialog line at x position
		lineRunes := []rune(line)
		dialogRunes := []rune(dialogLine)

		// Build new line with dialog overlay
		var newLine strings.Builder
		pos := 0

		// Characters before dialog
		for pos < x && pos < len(lineRunes) {
			newLine.WriteRune(lineRunes[pos])
			pos++
		}

		// Pad if needed
		for pos < x {
			newLine.WriteRune(' ')
			pos++
		}

		// Dialog content
		newLine.WriteString(string(dialogRunes))
		pos += lipgloss.Width(string(dialogRunes))

		// Characters after dialog
		if pos < len(lineRunes) {
			for i := pos; i < len(lineRunes); i++ {
				newLine.WriteRune(lineRunes[i])
			}
		}

		lines[lineIdx] = newLine.String()
	}

	// Add dim background effect
	t := theme.Current
	dimStyle := lipgloss.NewStyle().Foreground(t.Muted)

	for i := range lines {
		// Don't dim dialog lines
		if i < y || i >= y+len(dialogLines) {
			lines[i] = dimStyle.Render(lines[i])
		}
	}

	return strings.Join(lines[:m.height], "\n")
}
