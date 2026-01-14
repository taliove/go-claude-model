package ui

import (
	"errors"

	"ccm/internal/config"
	"ccm/internal/provider"
	"ccm/internal/ui/app"
	"ccm/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

// Errors
var (
	ErrCanceled = errors.New("operation canceled by user")
)

// TUIResult contains the result of running the TUI
type TUIResult struct {
	RunProvider string // Provider to run after TUI exits (empty if none)
}

// RunTUI launches the full-screen TUI and returns the result
func RunTUI() (*TUIResult, error) {
	m := app.NewApp()
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(app.AppModel)
	return &TUIResult{
		RunProvider: result.GetRunCommand(),
	}, nil
}

// ProviderItem represents a provider in the selection list
type ProviderItem struct {
	Name         string
	DisplayName  string
	IsConfigured bool
	IsDefault    bool
}

// providerSelectModel wraps ProviderListModel to implement tea.Model
type providerSelectModel struct {
	list components.ProviderListModel
}

func (m providerSelectModel) Init() tea.Cmd {
	return m.list.Init()
}

func (m providerSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m providerSelectModel) View() string {
	return m.list.View()
}

// SelectProvider shows an interactive provider selection with search
func SelectProvider(items []ProviderItem, label string) (string, error) {
	componentItems := make([]components.ProviderItem, len(items))
	for i, item := range items {
		componentItems[i] = components.ProviderItem{
			Name:         item.Name,
			DisplayName:  item.DisplayName,
			IsConfigured: item.IsConfigured,
			IsDefault:    item.IsDefault,
		}
	}

	innerModel := components.NewProviderList(componentItems, label)
	m := providerSelectModel{list: innerModel}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(providerSelectModel)
	if result.list.Canceled() {
		return "", ErrCanceled
	}
	if result.list.Selected() == nil {
		return "", ErrCanceled
	}
	return result.list.Selected().Name, nil
}

// PromptAPIKey shows a masked input for API key entry
func PromptAPIKey(label string) (string, error) {
	m := components.NewAPIKeyInput(label)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(components.TextInputModel)
	if result.Canceled() {
		return "", ErrCanceled
	}
	return result.Value(), nil
}

// PromptConfirm shows a Y/N confirmation dialog
func PromptConfirm(label string) bool {
	m := components.NewConfirm(label)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return false
	}

	result := finalModel.(components.ConfirmModel)
	return result.Confirmed()
}

// SelectAction shows an action menu and returns the selected index
func SelectAction(actions []string, label string) (int, error) {
	m := components.NewMenu(actions, label)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return -1, err
	}

	result := finalModel.(components.MenuModel)
	if result.Canceled() {
		return -1, ErrCanceled
	}
	return result.Selected(), nil
}

// BuildProviderItems builds provider list items from config
func BuildProviderItems(cfg *config.Config, includeUnconfigured bool) []ProviderItem {
	items := []ProviderItem{}

	// Add presets in order
	for _, name := range provider.PresetOrder {
		p := provider.Presets[name]
		isConfigured := false
		isDefault := false

		// Check if configured
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

		if includeUnconfigured || isConfigured {
			items = append(items, ProviderItem{
				Name:         name,
				DisplayName:  p.DisplayName,
				IsConfigured: isConfigured,
				IsDefault:    isDefault,
			})
		}
	}

	// Add custom providers
	for name, p := range cfg.Providers {
		if _, exists := provider.Presets[name]; !exists {
			isDefault := cfg.Default == name
			items = append(items, ProviderItem{
				Name:         name,
				DisplayName:  p.DisplayName,
				IsConfigured: true,
				IsDefault:    isDefault,
			})
		}
	}

	return items
}

// BuildConfiguredProviderItems builds only configured provider items
func BuildConfiguredProviderItems(cfg *config.Config) []ProviderItem {
	return BuildProviderItems(cfg, false)
}
