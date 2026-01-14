package components

import (
	"fmt"
	"strings"
	"time"

	"ccm/internal/provider"
	"ccm/internal/ui/messages"
	"ccm/internal/ui/theme"

	"github.com/charmbracelet/lipgloss"
)

// DetailPanelModel shows details of the selected provider
type DetailPanelModel struct {
	provider   *provider.Provider
	status     messages.ConnectionStatus
	latency    time.Duration
	statusText string
	width      int
}

// NewDetailPanel creates a new detail panel
func NewDetailPanel() DetailPanelModel {
	return DetailPanelModel{}
}

// SetProvider updates the displayed provider
func (m *DetailPanelModel) SetProvider(p *provider.Provider) {
	m.provider = p
}

// SetConnectionStatus updates the connection status
func (m *DetailPanelModel) SetConnectionStatus(status messages.ConnectionStatus, latency time.Duration, err error) {
	m.status = status
	m.latency = latency
	if err != nil {
		m.statusText = err.Error()
	} else {
		m.statusText = ""
	}
}

// SetWidth updates the panel width
func (m *DetailPanelModel) SetWidth(width int) {
	m.width = width
}

// View renders the detail panel
func (m DetailPanelModel) View() string {
	styles := theme.GetStyles()
	t := theme.Current

	if m.provider == nil {
		return styles.Muted.Render("No provider selected")
	}

	var b strings.Builder

	// Separator
	separator := strings.Repeat("─", m.width-2)
	b.WriteString(styles.Muted.Render(separator))
	b.WriteString("\n")

	// Model
	labelStyle := lipgloss.NewStyle().
		Foreground(t.Muted).
		Width(8)
	valueStyle := lipgloss.NewStyle().
		Foreground(t.Warning)

	b.WriteString(labelStyle.Render("Model:"))
	b.WriteString(" ")
	b.WriteString(valueStyle.Render(m.provider.Model))
	b.WriteString("\n")

	// URL
	b.WriteString(labelStyle.Render("URL:"))
	b.WriteString(" ")
	urlStyle := lipgloss.NewStyle().Foreground(t.Accent)
	b.WriteString(urlStyle.Render(m.provider.BaseURL))
	b.WriteString("\n")

	// Connection Status
	b.WriteString(labelStyle.Render("Status:"))
	b.WriteString(" ")

	switch m.status {
	case messages.ConnectionUnknown:
		b.WriteString(styles.Muted.Render("○ Not tested"))
	case messages.ConnectionTesting:
		b.WriteString(styles.Muted.Render("⟳ Testing..."))
	case messages.ConnectionOK:
		statusStyle := lipgloss.NewStyle().Foreground(t.Success)
		b.WriteString(statusStyle.Render(fmt.Sprintf("● Connected (%dms)", m.latency.Milliseconds())))
	case messages.ConnectionError:
		statusStyle := lipgloss.NewStyle().Foreground(t.Error)
		errText := "Connection failed"
		if m.statusText != "" {
			errText = m.statusText
		}
		b.WriteString(statusStyle.Render("● " + errText))
	}
	b.WriteString("\n")

	return b.String()
}
