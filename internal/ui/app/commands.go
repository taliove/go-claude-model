package app

import (
	"net/http"
	"time"

	"ccm/internal/config"
	"ccm/internal/ui/messages"

	tea "github.com/charmbracelet/bubbletea"
)

// testConnection tests the connection to a provider
func testConnection(name string) tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.Load()
		if err != nil {
			return messages.ConnectionResultMsg{
				Name:   name,
				Status: messages.ConnectionError,
				Error:  err,
			}
		}

		p, exists := cfg.Providers[name]
		if !exists {
			return messages.ConnectionResultMsg{
				Name:   name,
				Status: messages.ConnectionError,
				Error:  err,
			}
		}

		// Get effective API key
		apiKey := config.GetEffectiveAPIKey(name)
		if apiKey == "" {
			return messages.ConnectionResultMsg{
				Name:   name,
				Status: messages.ConnectionError,
				Error:  err,
			}
		}

		// Create HTTP client with timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Create request
		req, err := http.NewRequest("GET", p.BaseURL, nil)
		if err != nil {
			return messages.ConnectionResultMsg{
				Name:   name,
				Status: messages.ConnectionError,
				Error:  err,
			}
		}

		// Add authorization header
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		// Execute request and measure latency
		start := time.Now()
		resp, err := client.Do(req)
		latency := time.Since(start)

		if err != nil {
			return messages.ConnectionResultMsg{
				Name:    name,
				Status:  messages.ConnectionError,
				Latency: latency,
				Error:   err,
			}
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode >= 400 && resp.StatusCode != 404 {
			// 404 is often returned by API endpoints that don't support GET
			// We still consider it a successful connection
			return messages.ConnectionResultMsg{
				Name:    name,
				Status:  messages.ConnectionError,
				Latency: latency,
			}
		}

		return messages.ConnectionResultMsg{
			Name:    name,
			Status:  messages.ConnectionOK,
			Latency: latency,
		}
	}
}
