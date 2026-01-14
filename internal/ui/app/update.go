package app

import (
	"ccm/internal/config"
	"ccm/internal/provider"
	"ccm/internal/ui/dialogs"
	"ccm/internal/ui/messages"
	"ccm/internal/ui/theme"

	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Update component sizes
		m.header.SetWidth(msg.Width)
		m.providerList.SetSize(msg.Width, msg.Height-8) // header + detail + status
		m.detailPanel.SetWidth(msg.Width)
		m.statusBar.SetWidth(msg.Width)

		return m, nil

	case tea.KeyMsg:
		// Handle dialog input first
		if m.activeDialog != nil {
			return m.updateDialog(msg)
		}

		// Handle search mode
		if m.providerList.Searching() {
			return m.updateSearch(msg)
		}

		// Handle main navigation
		return m.updateMain(msg)

	case messages.ConnectionResultMsg:
		m.providerList.UpdateConnectionStatus(msg.Name, msg.Status, msg.Latency)

		// Update detail panel if this is the selected provider
		if selected := m.providerList.Selected(); selected != nil && selected.Name == msg.Name {
			m.detailPanel.SetConnectionStatus(msg.Status, msg.Latency, msg.Error)
		}

		return m, nil

	case messages.CloseDialogMsg:
		return m.handleDialogClose(msg)

	case messages.ThemeChangedMsg:
		m.statusBar.SetThemeIcon(msg.IsDark)
		return m, nil

	case messages.ConfigReloadedMsg:
		if msg.Error == nil {
			m.statusBar.SetMessage("Configuration reloaded", false)
		} else {
			m.statusBar.SetMessage("Failed to reload config", true)
		}
		return m, nil

	case messages.StatusMsg:
		m.statusBar.SetMessage(msg.Text, msg.IsError)
		return m, nil

	case providerSavedMsg:
		m.config = msg.config
		m.providers = buildProviderItems(msg.config)
		m.providerList.SetItems(m.providers)
		m.statusBar.SetMessage("Provider saved: "+msg.name, false)
		m.updateDetailPanel()
		return m, nil

	case providerRemovedMsg:
		m.config = msg.config
		m.providers = buildProviderItems(msg.config)
		m.providerList.SetItems(m.providers)
		m.statusBar.SetMessage("Provider removed: "+msg.name, false)
		m.updateDetailPanel()
		return m, nil

	case defaultSetMsg:
		m.config = msg.config
		m.providers = buildProviderItems(msg.config)
		m.providerList.SetItems(m.providers)
		m.statusBar.SetDefaultProvider(msg.name)
		m.statusBar.SetMessage("Default set: "+msg.name, false)
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m AppModel) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "?":
		// Show help dialog
		m.activeDialog = dialogs.NewHelpDialog()
		m.dialogType = messages.DialogHelp
		return m, nil

	case "/":
		// Start search
		var cmd tea.Cmd
		m.providerList, cmd = m.providerList.Update(messages.StartSearchMsg{})
		return m, cmd

	case "ctrl+t":
		// Toggle theme
		newTheme := theme.Toggle()
		m.statusBar.SetThemeIcon(newTheme.IsDark)
		return m, nil

	case "up", "k", "down", "j", "home", "g", "end", "G":
		// Navigation - update list and detail panel
		var cmd tea.Cmd
		m.providerList, cmd = m.providerList.Update(msg)
		m.updateDetailPanel()
		return m, cmd

	case "enter":
		// Run selected provider
		selected := m.providerList.Selected()
		if selected == nil {
			return m, nil
		}
		if !selected.IsConfigured {
			m.statusBar.SetMessage("Provider not configured. Press 'a' to add.", true)
			return m, nil
		}
		m.runCommand = selected.Name
		m.quitting = true
		return m, tea.Quit

	case "e":
		// Edit provider
		selected := m.providerList.Selected()
		if selected == nil {
			return m, nil
		}
		return m.openEditDialog(selected.Name)

	case "a":
		// Add/configure provider (same as edit for unconfigured)
		selected := m.providerList.Selected()
		if selected == nil {
			return m, nil
		}
		return m.openEditDialog(selected.Name)

	case "d":
		// Set as default
		selected := m.providerList.Selected()
		if selected == nil {
			return m, nil
		}
		if !selected.IsConfigured {
			m.statusBar.SetMessage("Cannot set unconfigured provider as default", true)
			return m, nil
		}
		return m, m.setDefault(selected.Name)

	case "t":
		// Test connection
		selected := m.providerList.Selected()
		if selected == nil {
			return m, nil
		}
		if !selected.IsConfigured {
			m.statusBar.SetMessage("Cannot test unconfigured provider", true)
			return m, nil
		}
		m.providerList.UpdateConnectionStatus(selected.Name, messages.ConnectionTesting, 0)
		m.detailPanel.SetConnectionStatus(messages.ConnectionTesting, 0, nil)
		return m, testConnection(selected.Name)

	case "r":
		// Remove provider
		selected := m.providerList.Selected()
		if selected == nil {
			return m, nil
		}
		if _, exists := m.config.Providers[selected.Name]; !exists {
			m.statusBar.SetMessage("Provider not configured", true)
			return m, nil
		}
		return m.openRemoveDialog(selected.Name, selected.DisplayName)
	}

	return m, nil
}

func (m AppModel) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.providerList, cmd = m.providerList.Update(msg)

	// Update detail panel when search results change
	m.updateDetailPanel()

	return m, cmd
}

func (m AppModel) updateDialog(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	updatedDialog, cmd := m.activeDialog.Update(msg)
	m.activeDialog = updatedDialog.(dialogs.Dialog)
	return m, cmd
}

func (m AppModel) handleDialogClose(msg messages.CloseDialogMsg) (tea.Model, tea.Cmd) {
	dialogType := m.dialogType
	m.activeDialog = nil
	m.dialogType = messages.DialogNone

	switch dialogType {
	case messages.DialogEdit:
		if msg.Result != nil {
			if p, ok := msg.Result.(provider.Provider); ok {
				return m, m.saveProvider(p)
			}
		}

	case messages.DialogConfirm:
		if confirmed, ok := msg.Result.(bool); ok && confirmed {
			// Get the provider name from status message context
			if selected := m.providerList.Selected(); selected != nil {
				return m, m.removeProvider(selected.Name)
			}
		}
	}

	return m, nil
}

func (m *AppModel) updateDetailPanel() {
	selected := m.providerList.Selected()
	if selected == nil {
		m.detailPanel.SetProvider(nil)
		return
	}

	// Try to get from config first
	if p, exists := m.config.Providers[selected.Name]; exists {
		m.detailPanel.SetProvider(&p)
		return
	}

	// Fall back to preset
	if preset, exists := provider.Presets[selected.Name]; exists {
		m.detailPanel.SetProvider(&preset)
		return
	}

	m.detailPanel.SetProvider(nil)
}

func (m AppModel) openEditDialog(name string) (tea.Model, tea.Cmd) {
	var p provider.Provider

	// Get existing config or preset
	if cp, exists := m.config.Providers[name]; exists {
		p = cp
	} else if preset, exists := provider.Presets[name]; exists {
		p = preset
	} else {
		return m, nil
	}

	m.activeDialog = dialogs.NewEditDialog(p)
	m.dialogType = messages.DialogEdit
	return m, nil
}

func (m AppModel) openRemoveDialog(name, displayName string) (tea.Model, tea.Cmd) {
	m.activeDialog = dialogs.RemoveConfirmDialog(name, displayName)
	m.dialogType = messages.DialogConfirm
	return m, nil
}

func (m AppModel) saveProvider(p provider.Provider) tea.Cmd {
	return func() tea.Msg {
		if err := config.AddProvider(p); err != nil {
			return messages.StatusMsg{Text: "Failed to save provider", IsError: true}
		}

		// Reload config
		cfg, err := config.Load()
		if err != nil {
			return messages.StatusMsg{Text: "Failed to reload config", IsError: true}
		}

		return providerSavedMsg{config: cfg, name: p.Name}
	}
}

func (m AppModel) removeProvider(name string) tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.Load()
		if err != nil {
			return messages.StatusMsg{Text: "Failed to load config", IsError: true}
		}

		delete(cfg.Providers, name)

		if err := config.Save(cfg); err != nil {
			return messages.StatusMsg{Text: "Failed to save config", IsError: true}
		}

		return providerRemovedMsg{config: cfg, name: name}
	}
}

func (m AppModel) setDefault(name string) tea.Cmd {
	return func() tea.Msg {
		if err := config.SetDefault(name); err != nil {
			return messages.StatusMsg{Text: "Failed to set default", IsError: true}
		}

		cfg, err := config.Load()
		if err != nil {
			return messages.StatusMsg{Text: "Failed to reload config", IsError: true}
		}

		return defaultSetMsg{config: cfg, name: name}
	}
}

// Internal messages for config updates
type providerSavedMsg struct {
	config *config.Config
	name   string
}

type providerRemovedMsg struct {
	config *config.Config
	name   string
}

type defaultSetMsg struct {
	config *config.Config
	name   string
}
