package ui

import (
	"fmt"
	"strings"

	"ccm/internal/config"
	"ccm/internal/provider"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// ProviderItem 供应商列表项
type ProviderItem struct {
	Name         string
	DisplayName  string
	IsConfigured bool
	IsDefault    bool
}

// SelectProvider 使用箭头键选择供应商
func SelectProvider(items []ProviderItem, label string) (string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Name | cyan }} ({{ .DisplayName }}){{ if .IsConfigured }} {{ \"✓\" | green }}{{ end }}{{ if .IsDefault }} {{ \"★\" | yellow }}{{ end }}",
		Inactive: "  {{ .Name }} ({{ .DisplayName }}){{ if .IsConfigured }} {{ \"✓\" | green }}{{ end }}{{ if .IsDefault }} {{ \"★\" | yellow }}{{ end }}",
		Selected: "{{ \"✓\" | green }} {{ .Name | cyan }} ({{ .DisplayName }})",
		Details: `
--------- 供应商详情 ----------
{{ "名称:" | faint }}	{{ .Name }}
{{ "显示名:" | faint }}	{{ .DisplayName }}
{{ "状态:" | faint }}	{{ if .IsConfigured }}已配置{{ else }}未配置{{ end }}`,
	}

	searcher := func(input string, index int) bool {
		item := items[index]
		name := strings.ToLower(item.Name)
		displayName := strings.ToLower(item.DisplayName)
		input = strings.ToLower(input)
		return strings.Contains(name, input) || strings.Contains(displayName, input)
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return items[i].Name, nil
}

// PromptAPIKey 带掩码的 API Key 输入
func PromptAPIKey(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
		Validate: func(input string) error {
			if len(input) < 1 {
				return fmt.Errorf("API Key 不能为空")
			}
			return nil
		},
	}

	return prompt.Run()
}

// PromptConfirm Y/N 确认
func PromptConfirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	return strings.ToLower(result) == "y"
}

// SelectAction 选择操作菜单
func SelectAction(actions []string, label string) (int, error) {
	prompt := promptui.Select{
		Label: label,
		Items: actions,
		Size:  6,
	}

	i, _, err := prompt.Run()
	return i, err
}

// BuildProviderItems 构建供应商列表项
func BuildProviderItems(cfg *config.Config, includeUnconfigured bool) []ProviderItem {
	items := []ProviderItem{}

	// 按预置顺序添加
	for _, name := range provider.PresetOrder {
		p := provider.Presets[name]
		isConfigured := false
		isDefault := false

		// 检查是否已配置
		if cp, exists := cfg.Providers[name]; exists && cp.APIKey != "" {
			isConfigured = true
		}

		// 检查环境变量
		if config.GetEnvAPIKey(name) != "" {
			isConfigured = true
		}

		// 检查是否是默认供应商
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

	// 添加用户自定义供应商
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

// BuildConfiguredProviderItems 只构建已配置的供应商列表
func BuildConfiguredProviderItems(cfg *config.Config) []ProviderItem {
	return BuildProviderItems(cfg, false)
}

// 颜色函数供模板使用
func init() {
	promptui.FuncMap["cyan"] = func(s string) string {
		return color.CyanString(s)
	}
	promptui.FuncMap["green"] = func(s string) string {
		return color.GreenString(s)
	}
	promptui.FuncMap["yellow"] = func(s string) string {
		return color.YellowString(s)
	}
	promptui.FuncMap["red"] = func(s string) string {
		return color.RedString(s)
	}
}
