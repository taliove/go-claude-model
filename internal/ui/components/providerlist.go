package components

import (
	"fmt"
	"strings"
	"time"

	"ccm/internal/ui/messages"
	"ccm/internal/ui/theme"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProviderItem is an alias for backward compatibility with ui.go
type ProviderItem = ProviderListItem

// ProviderListItem represents a provider in the list
type ProviderListItem struct {
	Name         string
	DisplayName  string
	IsConfigured bool
	IsDefault    bool
	Status       messages.ConnectionStatus
	Latency      time.Duration
}

// ProviderListModel is the provider selection list with search
type ProviderListModel struct {
	label    string
	items    []ProviderListItem
	filtered []ProviderListItem
	cursor   int
	offset   int // scroll offset

	// Search
	searchInput textinput.Model
	searching   bool
	canceled    bool
	quitting    bool

	// Dimensions
	width  int
	height int
}

// NewProviderList creates a new provider list (with label for backward compatibility)
func NewProviderList(items []ProviderListItem, label string) ProviderListModel {
	ti := textinput.New()
	ti.Placeholder = "Search providers..."
	ti.CharLimit = 50

	return ProviderListModel{
		label:       label,
		items:       items,
		filtered:    items,
		cursor:      0,
		searchInput: ti,
	}
}

// NewProviderListSimple creates a provider list without label (for TUI app)
func NewProviderListSimple(items []ProviderListItem) ProviderListModel {
	return NewProviderList(items, "")
}

// SetItems updates the provider list
func (m *ProviderListModel) SetItems(items []ProviderListItem) {
	m.items = items
	m.filterItems()
}

// SetSize updates the list dimensions
func (m *ProviderListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Selected returns the currently selected item
func (m ProviderListModel) Selected() *ProviderListItem {
	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
		return nil
	}
	item := m.filtered[m.cursor]
	return &item
}

// Searching returns whether search is active
func (m ProviderListModel) Searching() bool {
	return m.searching
}

// Canceled returns whether the user canceled the selection
func (m ProviderListModel) Canceled() bool {
	return m.canceled
}

// Init implements tea.Model
func (m ProviderListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ProviderListModel) Update(msg tea.Msg) (ProviderListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			return m.updateSearch(msg)
		}
		return m.updateNavigation(msg)
	case messages.StartSearchMsg:
		m.searching = true
		m.searchInput.Focus()
		return m, textinput.Blink
	case messages.EndSearchMsg:
		m.searching = false
		m.searchInput.SetValue("")
		m.searchInput.Blur()
		m.filterItems()
		return m, nil
	}
	return m, nil
}

func (m ProviderListModel) updateSearch(msg tea.KeyMsg) (ProviderListModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searching = false
		m.searchInput.SetValue("")
		m.searchInput.Blur()
		m.filterItems()
		return m, nil
	case "enter":
		m.searching = false
		m.searchInput.Blur()
		// Select current item
		if len(m.filtered) > 0 {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	case "up", "ctrl+p":
		if m.cursor > 0 {
			m.cursor--
			m.ensureVisible()
		}
		return m, nil
	case "down", "ctrl+n":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
			m.ensureVisible()
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.filterItems()
	return m, cmd
}

func (m ProviderListModel) updateNavigation(msg tea.KeyMsg) (ProviderListModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.ensureVisible()
		}
	case "down", "j":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
			m.ensureVisible()
		}
	case "home", "g":
		m.cursor = 0
		m.offset = 0
	case "end", "G":
		m.cursor = len(m.filtered) - 1
		m.ensureVisible()
	case "enter":
		// Select current item (for standalone use)
		if len(m.filtered) > 0 {
			m.quitting = true
			return m, tea.Quit
		}
	case "esc", "ctrl+c":
		m.canceled = true
		m.quitting = true
		return m, tea.Quit
	case "/":
		m.searching = true
		m.searchInput.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m *ProviderListModel) filterItems() {
	query := strings.ToLower(m.searchInput.Value())
	if query == "" {
		m.filtered = m.items
	} else {
		m.filtered = []ProviderListItem{}
		for _, item := range m.items {
			name := strings.ToLower(item.Name)
			displayName := strings.ToLower(item.DisplayName)
			if strings.Contains(name, query) || strings.Contains(displayName, query) {
				m.filtered = append(m.filtered, item)
			}
		}
	}

	// Adjust cursor if out of bounds
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
	m.ensureVisible()
}

func (m *ProviderListModel) ensureVisible() {
	visibleHeight := m.height - 4 // Account for header, search, borders
	if visibleHeight < 1 {
		visibleHeight = 7
	}

	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+visibleHeight {
		m.offset = m.cursor - visibleHeight + 1
	}
}

// View renders the provider list
func (m ProviderListModel) View() string {
	if m.quitting {
		return ""
	}

	styles := theme.GetStyles()
	t := theme.Current

	var b strings.Builder

	// Label (if provided, for standalone mode)
	if m.label != "" {
		b.WriteString(styles.Title.Render(m.label))
		b.WriteString("\n")
	}

	// Title with count
	configured := 0
	for _, item := range m.items {
		if item.IsConfigured {
			configured++
		}
	}
	titleText := fmt.Sprintf("Providers (%d/%d configured)", configured, len(m.items))

	searchHint := ""
	if m.searching {
		searchHint = " [searching]"
	} else {
		searchHint = " [/]"
	}

	title := lipgloss.JoinHorizontal(
		lipgloss.Top,
		styles.Title.Render(titleText),
		styles.Muted.Render(searchHint),
	)
	b.WriteString(title)
	b.WriteString("\n")

	// Search input (if searching)
	if m.searching {
		searchStyle := lipgloss.NewStyle().
			Foreground(t.Primary).
			Padding(0, 1)
		b.WriteString(searchStyle.Render(m.searchInput.View()))
		b.WriteString("\n")
	}

	// Separator
	separator := strings.Repeat("─", m.width-2)
	b.WriteString(styles.Muted.Render(separator))
	b.WriteString("\n")

	// Calculate visible range
	visibleHeight := m.height - 5
	if m.searching {
		visibleHeight--
	}
	if visibleHeight < 1 {
		visibleHeight = 7
	}

	endIdx := m.offset + visibleHeight
	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	// Items
	for i := m.offset; i < endIdx; i++ {
		item := m.filtered[i]
		isSelected := i == m.cursor

		// Cursor
		cursor := "  "
		if isSelected {
			cursor = styles.Cursor.Render("▸ ")
		}

		// Name style
		nameStyle := styles.Normal
		if isSelected {
			nameStyle = styles.Selected
		}

		// Status indicator
		var status string
		if item.IsConfigured {
			status = styles.Success.Render(" ✓")
		} else {
			status = styles.Error.Render(" ✗")
		}

		// Default marker
		defaultMark := ""
		if item.IsDefault {
			defaultMark = styles.Warning.Render(" ★")
		}

		// Connection status
		connStatus := ""
		switch item.Status {
		case messages.ConnectionTesting:
			connStatus = styles.Muted.Render(" ⟳")
		case messages.ConnectionOK:
			connStatus = styles.Success.Render(fmt.Sprintf(" ●%dms", item.Latency.Milliseconds()))
		case messages.ConnectionError:
			connStatus = styles.Error.Render(" ●")
		}

		// Format line
		name := nameStyle.Render(fmt.Sprintf("%-12s", item.Name))
		displayName := styles.Muted.Render(item.DisplayName)

		line := fmt.Sprintf("%s%s %s%s%s%s",
			cursor,
			name,
			displayName,
			status,
			defaultMark,
			connStatus,
		)

		b.WriteString(line)
		b.WriteString("\n")
	}

	// Scroll indicators
	if m.offset > 0 {
		b.WriteString(styles.Muted.Render("  ↑ more above\n"))
	}
	if endIdx < len(m.filtered) {
		b.WriteString(styles.Muted.Render("  ↓ more below\n"))
	}

	return b.String()
}

// UpdateConnectionStatus updates the connection status for a provider
func (m *ProviderListModel) UpdateConnectionStatus(name string, status messages.ConnectionStatus, latency time.Duration) {
	for i := range m.items {
		if m.items[i].Name == name {
			m.items[i].Status = status
			m.items[i].Latency = latency
			break
		}
	}
	for i := range m.filtered {
		if m.filtered[i].Name == name {
			m.filtered[i].Status = status
			m.filtered[i].Latency = latency
			break
		}
	}
}
