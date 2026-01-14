package app

import (
	"ccm/internal/config"
	"ccm/internal/provider"
	"ccm/internal/ui/components"
	"ccm/internal/ui/dialogs"
	"ccm/internal/ui/messages"
	"ccm/internal/ui/theme"

	tea "github.com/charmbracelet/bubbletea"
)

// AppModel is the root TUI model
type AppModel struct {
	// Data
	config    *config.Config
	providers []components.ProviderListItem

	// Components
	header       components.HeaderModel
	providerList components.ProviderListModel
	detailPanel  components.DetailPanelModel
	statusBar    components.StatusBarModel

	// Dialog state
	activeDialog dialogs.Dialog
	dialogType   messages.DialogType

	// UI state
	width      int
	height     int
	ready      bool
	quitting   bool
	runCommand string // Provider to run after quit
}

// NewApp creates a new app model
func NewApp() AppModel {
	// Initialize theme
	theme.Set(theme.DetectSystemTheme())

	// Load config
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{
			Providers: make(map[string]provider.Provider),
		}
	}

	// Build provider items
	items := buildProviderItems(cfg)

	// Create components
	m := AppModel{
		config:       cfg,
		providers:    items,
		header:       components.NewHeader(),
		providerList: components.NewProviderListSimple(items),
		detailPanel:  components.NewDetailPanel(),
		statusBar:    components.NewStatusBar(),
	}

	// Set initial state
	m.statusBar.SetDefaultProvider(cfg.Default)
	m.statusBar.SetThemeIcon(theme.Current.IsDark)

	// Update detail panel with first provider
	if len(items) > 0 {
		if selected := m.providerList.Selected(); selected != nil {
			if p, exists := cfg.Providers[selected.Name]; exists {
				m.detailPanel.SetProvider(&p)
			} else if preset, exists := provider.Presets[selected.Name]; exists {
				m.detailPanel.SetProvider(&preset)
			}
		}
	}

	return m
}

// buildProviderItems converts config to list items
func buildProviderItems(cfg *config.Config) []components.ProviderListItem {
	items := []components.ProviderListItem{}

	// Add presets in order
	for _, name := range provider.PresetOrder {
		p := provider.Presets[name]
		isConfigured := false
		isDefault := false

		// Check if configured in config
		if cp, exists := cfg.Providers[name]; exists && cp.APIKey != "" {
			isConfigured = true
		}

		// Check environment variable
		if config.GetEnvAPIKey(name) != "" {
			isConfigured = true
		}

		// Check if default
		if cfg.Default == name {
			isDefault = true
		}

		items = append(items, components.ProviderListItem{
			Name:         name,
			DisplayName:  p.DisplayName,
			IsConfigured: isConfigured,
			IsDefault:    isDefault,
			Status:       messages.ConnectionUnknown,
		})
	}

	// Add custom providers (not in presets)
	for name, p := range cfg.Providers {
		if _, exists := provider.Presets[name]; !exists {
			isDefault := cfg.Default == name
			items = append(items, components.ProviderListItem{
				Name:         name,
				DisplayName:  p.DisplayName,
				IsConfigured: true,
				IsDefault:    isDefault,
				Status:       messages.ConnectionUnknown,
			})
		}
	}

	return items
}

// GetRunCommand returns the provider to run after quit
func (m AppModel) GetRunCommand() string {
	return m.runCommand
}

// Init implements tea.Model
func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
	)
}
