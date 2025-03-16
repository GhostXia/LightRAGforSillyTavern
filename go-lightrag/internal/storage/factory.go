package storage

import (
	"github.com/HerSophia/go-lightrag/internal/config"
	"github.com/HerSophia/go-lightrag/internal/storage/graph"
	"github.com/HerSophia/go-lightrag/internal/storage/kv"
	"github.com/HerSophia/go-lightrag/internal/storage/vector"
	"github.com/HerSophia/go-lightrag/pkg/interfaces"
	"go.uber.org/zap"
)

// NewVectorStorage 根据配置创建向量存储
func NewVectorStorage(cfg config.StorageConfig, logger *zap.Logger) (interfaces.VectorStorage, error) {
	switch cfg.VectorType {
	case "memory":
		logger.Info("初始化内存向量存储")
		return vector.NewMemoryVectorStorage(), nil
	// TODO: 添加其他向量存储实现，如postgres、milvus等
	default:
		logger.Warn("未知的向量存储类型，使用内存向量存储", zap.String("type", cfg.VectorType))
		return vector.NewMemoryVectorStorage(), nil
	}
}

// NewKVStorage 根据配置创建KV存储
func NewKVStorage(cfg config.StorageConfig, logger *zap.Logger) (interfaces.KVStorage, error) {
	switch cfg.KvType {
	case "memory":
		logger.Info("初始化内存KV存储")
		return kv.NewMemoryKVStorage(), nil
	// TODO: 添加其他KV存储实现，如postgres、mongo等
	default:
		logger.Warn("未知的KV存储类型，使用内存KV存储", zap.String("type", cfg.KvType))
		return kv.NewMemoryKVStorage(), nil
	}
}

// NewGraphStorage 根据配置创建图存储
func NewGraphStorage(cfg config.StorageConfig, logger *zap.Logger) (interfaces.GraphStorage, error) {
	switch cfg.GraphType {
	case "memory":
		logger.Info("初始化内存图存储")
		return graph.NewMemoryGraphStorage(), nil
	// TODO: 添加其他图存储实现，如neo4j、postgres等
	default:
		logger.Warn("未知的图存储类型，使用内存图存储", zap.String("type", cfg.GraphType))
		return graph.NewMemoryGraphStorage(), nil
	}
}
