package messages

import (
	"time"

	"ccm/internal/provider"
)

// Navigation messages
type (
	// CursorUpMsg moves cursor up
	CursorUpMsg struct{}

	// CursorDownMsg moves cursor down
	CursorDownMsg struct{}

	// CursorHomeMsg jumps to first item
	CursorHomeMsg struct{}

	// CursorEndMsg jumps to last item
	CursorEndMsg struct{}
)

// Action messages
type (
	// RunProviderMsg triggers running Claude with a provider
	RunProviderMsg struct {
		Name string
	}

	// EditProviderMsg opens edit dialog for a provider
	EditProviderMsg struct {
		Name string
	}

	// SetDefaultMsg sets a provider as default
	SetDefaultMsg struct {
		Name string
	}

	// TestConnectionMsg triggers connection test for a provider
	TestConnectionMsg struct {
		Name string
	}

	// RemoveProviderMsg triggers provider removal
	RemoveProviderMsg struct {
		Name string
	}

	// AddProviderMsg triggers adding a new provider
	AddProviderMsg struct {
		Name string
	}
)

// Connection status
type ConnectionStatus int

const (
	ConnectionUnknown ConnectionStatus = iota
	ConnectionTesting
	ConnectionOK
	ConnectionError
)

// ConnectionResultMsg contains the result of a connection test
type ConnectionResultMsg struct {
	Name    string
	Status  ConnectionStatus
	Latency time.Duration
	Error   error
}

// Dialog messages
type (
	// OpenDialogMsg opens a dialog
	OpenDialogMsg struct {
		Type    DialogType
		Context interface{}
	}

	// CloseDialogMsg closes the current dialog
	CloseDialogMsg struct {
		Result interface{}
	}

	// DialogType identifies the type of dialog
	DialogType int
)

const (
	DialogNone DialogType = iota
	DialogEdit
	DialogConfirm
	DialogHelp
)

// Search messages
type (
	// StartSearchMsg activates search mode
	StartSearchMsg struct{}

	// EndSearchMsg deactivates search mode
	EndSearchMsg struct{}

	// SearchQueryMsg updates search query
	SearchQueryMsg struct {
		Query string
	}
)

// Theme messages
type (
	// ThemeToggleMsg toggles between dark and light themes
	ThemeToggleMsg struct{}

	// ThemeChangedMsg indicates theme has changed
	ThemeChangedMsg struct {
		IsDark bool
	}
)

// Provider data messages
type (
	// ProviderUpdatedMsg indicates a provider was updated
	ProviderUpdatedMsg struct {
		Provider provider.Provider
	}

	// ProviderRemovedMsg indicates a provider was removed
	ProviderRemovedMsg struct {
		Name string
	}

	// ConfigReloadMsg requests config reload
	ConfigReloadMsg struct{}

	// ConfigReloadedMsg indicates config was reloaded
	ConfigReloadedMsg struct {
		Error error
	}
)

// App lifecycle messages
type (
	// QuitMsg triggers app quit
	QuitMsg struct{}

	// ErrorMsg reports an error
	ErrorMsg struct {
		Error error
	}

	// StatusMsg shows a temporary status message
	StatusMsg struct {
		Text    string
		IsError bool
	}
)
