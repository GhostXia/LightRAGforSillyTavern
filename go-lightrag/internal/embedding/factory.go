package embedding

import (
	"fmt"

	"github.com/HerSophia/go-lightrag/internal/config"
	"github.com/HerSophia/go-lightrag/pkg/interfaces"
	"go.uber.org/zap"
)

// NewEmbeddingClient 根据配置创建嵌入模型客户端
func NewEmbeddingClient(cfg config.EmbeddingConfig, logger *zap.Logger) (interfaces.EmbeddingClient, error) {
	logger.Info("初始化嵌入模型客户端",
		zap.String("provider", cfg.Provider),
		zap.String("model", cfg.Model))

	switch cfg.Provider {
	case "openai":
		return NewOpenAIClient(cfg.ApiKey, cfg.BaseUrl, cfg.Model)
	// TODO: 添加其他嵌入模型提供商的支持，如ollama等
	default:
		return nil, fmt.Errorf("不支持的嵌入模型提供商: %s", cfg.Provider)
	}
}
