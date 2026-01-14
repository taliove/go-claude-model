package dialogs

import tea "github.com/charmbracelet/bubbletea"

// Dialog is the interface for all dialog types
type Dialog interface {
	tea.Model
	Title() string
	Width() int
	Height() int
	View() string
}
