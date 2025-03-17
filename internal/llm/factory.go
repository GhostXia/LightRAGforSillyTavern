package llm

import (
	"fmt"

	"github.com/HerSophia/go-lightrag/internal/config"
	"github.com/HerSophia/go-lightrag/pkg/interfaces"
	"go.uber.org/zap"
)

// NewLLMClient 根据配置创建LLM客户端
func NewLLMClient(cfg config.LlmConfig, logger *zap.Logger) (interfaces.LLMClient, error) {
	logger.Info("初始化LLM客户端",
		zap.String("provider", cfg.Provider),
		zap.String("model", cfg.Model))

	switch cfg.Provider {
	case "openai":
		return NewOpenAIClient(cfg.ApiKey, cfg.BaseUrl, cfg.Model)
	// TODO: 添加其他LLM提供商的支持，如ollama等
	default:
		return nil, fmt.Errorf("不支持的LLM提供商: %s", cfg.Provider)
	}
}
