package provider

// ProviderType 供应商类型
type ProviderType string

const (
	// TypeNativeModel 自有模型 - 供应商提供自己的模型
	TypeNativeModel ProviderType = "native"
	// TypeProxy 代理服务 - 代理 Anthropic API
	TypeProxy ProviderType = "proxy"
)

// Provider 供应商配置
type Provider struct {
	Name        string       `yaml:"name"`         // 供应商名称（用于命令行）
	DisplayName string       `yaml:"display_name"` // 显示名称（中文）
	APIKey      string       `yaml:"api_key"`      // API 密钥
	BaseURL     string       `yaml:"base_url"`     // API 基础 URL
	Model       string       `yaml:"model"`        // 默认模型
	KeyURL      string       `yaml:"key_url"`      // 获取 API Key 的网址
	Type        ProviderType `yaml:"type"`         // 供应商类型
}

// 预置供应商列表（用户只需填 API Key）
var Presets = map[string]Provider{
	"doubao": {
		Name:        "doubao",
		DisplayName: "豆包（字节跳动）",
		BaseURL:     "https://ark.cn-beijing.volces.com/api/compatible",
		Model:       "doubao-seed-code-preview-latest",
		KeyURL:      "https://console.volcengine.com/ark",
		Type:        TypeNativeModel,
	},
	"deepseek": {
		Name:        "deepseek",
		DisplayName: "DeepSeek（深度求索）",
		BaseURL:     "https://api.deepseek.com",
		Model:       "deepseek-chat",
		KeyURL:      "https://platform.deepseek.com",
		Type:        TypeNativeModel,
	},
	"qwen": {
		Name:        "qwen",
		DisplayName: "通义千问（阿里云）",
		BaseURL:     "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:       "qwen-max",
		KeyURL:      "https://dashscope.console.aliyun.com",
		Type:        TypeNativeModel,
	},
	"kimi": {
		Name:        "kimi",
		DisplayName: "Kimi（月之暗面）",
		BaseURL:     "https://api.moonshot.cn/v1",
		Model:       "moonshot-v1-auto",
		KeyURL:      "https://platform.moonshot.cn",
		Type:        TypeNativeModel,
	},
	"siliconflow": {
		Name:        "siliconflow",
		DisplayName: "硅基流动",
		BaseURL:     "https://api.siliconflow.cn/v1",
		Model:       "deepseek-ai/DeepSeek-V3",
		KeyURL:      "https://cloud.siliconflow.cn",
		Type:        TypeNativeModel,
	},
	"glm": {
		Name:        "glm",
		DisplayName: "智谱GLM",
		BaseURL:     "https://open.bigmodel.cn/api/paas/v4",
		Model:       "glm-4-plus",
		KeyURL:      "https://open.bigmodel.cn",
		Type:        TypeNativeModel,
	},
	"wanjie": {
		Name:        "wanjie",
		DisplayName: "万界（Claude 代理）",
		BaseURL:     "https://maas-openapi.wanjiedata.com/api/anthropic",
		Model:       "claude-opus-4-5-20251101",
		KeyURL:      "https://maas-openapi.wanjiedata.com",
		Type:        TypeProxy,
	},
}

// PresetOrder 预置供应商的显示顺序
var PresetOrder = []string{"doubao", "deepseek", "qwen", "kimi", "siliconflow", "glm", "wanjie"}
