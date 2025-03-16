package lightrag

import (
	"context"
	"errors"

	"github.com/HerSophia/go-lightrag/internal/config"
	"github.com/HerSophia/go-lightrag/internal/embedding"
	"github.com/HerSophia/go-lightrag/internal/llm"
	"github.com/HerSophia/go-lightrag/internal/storage/graph"
	"github.com/HerSophia/go-lightrag/internal/storage/kv"
	"github.com/HerSophia/go-lightrag/internal/storage/vector"
	"github.com/HerSophia/go-lightrag/pkg/interfaces"
	"go.uber.org/zap"
)

// LightRAG 是RAG系统的核心结构
type LightRAG struct {
	workingDir     string
	maxTokens      int
	maxEmbedTokens int
	chunkSize      int
	chunkOverlap   int
	topK           int
	logger         *zap.Logger

	// 存储组件
	vectorStorage interfaces.VectorStorage
	kvStorage     interfaces.KVStorage
	graphStorage  interfaces.GraphStorage

	// 模型组件
	llmClient   interfaces.LLMClient
	embedClient interfaces.EmbeddingClient
}

// NewLightRAG 创建新的LightRAG实例
func NewLightRAG(workingDir string, maxTokens, maxEmbedTokens, chunkSize, chunkOverlap, topK int, logger *zap.Logger) *LightRAG {
	return &LightRAG{
		workingDir:     workingDir,
		maxTokens:      maxTokens,
		maxEmbedTokens: maxEmbedTokens,
		chunkSize:      chunkSize,
		chunkOverlap:   chunkOverlap,
		topK:           topK,
		logger:         logger,
	}
}

// InitStorage 初始化存储组件
func (r *LightRAG) InitStorage(cfg config.StorageConfig) error {
	r.logger.Info("初始化存储组件")

	// 初始化向量存储
	switch cfg.VectorType {
	case "memory":
		r.vectorStorage = vector.NewMemoryVectorStorage()
		r.logger.Info("使用内存向量存储")
	// TODO: 添加其他向量存储实现，如postgres、milvus等
	default:
		r.logger.Warn("未知的向量存储类型，使用内存向量存储", zap.String("type", cfg.VectorType))
		r.vectorStorage = vector.NewMemoryVectorStorage()
	}

	// 初始化KV存储
	switch cfg.KvType {
	case "memory":
		r.kvStorage = kv.NewMemoryKVStorage()
		r.logger.Info("使用内存KV存储")
	// TODO: 添加其他KV存储实现，如postgres、mongo等
	default:
		r.logger.Warn("未知的KV存储类型，使用内存KV存储", zap.String("type", cfg.KvType))
		r.kvStorage = kv.NewMemoryKVStorage()
	}

	// 初始化图存储
	switch cfg.GraphType {
	case "memory":
		r.graphStorage = graph.NewMemoryGraphStorage()
		r.logger.Info("使用内存图存储")
	// TODO: 添加其他图存储实现，如neo4j、postgres等
	default:
		r.logger.Warn("未知的图存储类型，使用内存图存储", zap.String("type", cfg.GraphType))
		r.graphStorage = graph.NewMemoryGraphStorage()
	}

	return nil
}

// InitLLM 初始化LLM客户端
func (r *LightRAG) InitLLM(cfg config.LlmConfig) error {
	r.logger.Info("初始化LLM客户端")

	// 使用LLM工厂函数创建客户端
	client, err := llm.NewLLMClient(cfg, r.logger)
	if err != nil {
		r.logger.Error("初始化LLM客户端失败", zap.Error(err))
		return err
	}

	r.llmClient = client
	r.logger.Info("LLM客户端初始化成功", zap.String("provider", cfg.Provider), zap.String("model", cfg.Model))
	return nil
}

// InitEmbedding 初始化嵌入模型客户端
func (r *LightRAG) InitEmbedding(cfg config.EmbeddingConfig) error {
	r.logger.Info("初始化嵌入模型客户端")

	// 使用嵌入模型工厂函数创建客户端
	client, err := embedding.NewEmbeddingClient(cfg, r.logger)
	if err != nil {
		r.logger.Error("初始化嵌入模型客户端失败", zap.Error(err))
		return err
	}

	r.embedClient = client
	r.logger.Info("嵌入模型客户端初始化成功", zap.String("provider", cfg.Provider), zap.String("model", cfg.Model))
	return nil
}

// Insert 插入文本到RAG系统
func (r *LightRAG) Insert(text string) error {
	// TODO: 实现文本插入逻辑
	return nil
}

// Embed 生成文本嵌入
func (r *LightRAG) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if r.embedClient == nil {
		return nil, errors.New("嵌入模型客户端未初始化")
	}

	return r.embedClient.Embed(ctx, texts)
}

// CompleteStream 流式完成文本生成
func (r *LightRAG) CompleteStream(ctx context.Context, prompt string, systemPrompt string, maxTokens int, callback func(string) error) error {
	if r.llmClient == nil {
		return errors.New("LLM客户端未初始化")
	}

	return r.llmClient.CompleteStream(ctx, prompt, systemPrompt, maxTokens, callback)
}

// Query 查询RAG系统
func (r *LightRAG) Query(query string, params QueryParams) (string, error) {
	if r.llmClient == nil {
		return "", errors.New("LLM客户端未初始化")
	}

	// 简单实现，直接调用LLM，不使用RAG
	// TODO: 实现完整的RAG查询逻辑
	return r.llmClient.Complete(context.Background(), query, "", params.MaxTokenForGlobal)
}

// Close 关闭RAG系统及其资源
func (r *LightRAG) Close() error {
	r.logger.Info("关闭RAG系统资源")

	// 关闭存储资源
	if r.vectorStorage != nil {
		if err := r.vectorStorage.Close(); err != nil {
			r.logger.Error("关闭向量存储失败", zap.Error(err))
		}
	}

	if r.kvStorage != nil {
		if err := r.kvStorage.Close(); err != nil {
			r.logger.Error("关闭KV存储失败", zap.Error(err))
		}
	}

	if r.graphStorage != nil {
		if err := r.graphStorage.Close(); err != nil {
			r.logger.Error("关闭图存储失败", zap.Error(err))
		}
	}

	// 关闭模型资源
	if r.llmClient != nil {
		if err := r.llmClient.Close(); err != nil {
			r.logger.Error("关闭LLM客户端失败", zap.Error(err))
		}
	}

	if r.embedClient != nil {
		if err := r.embedClient.Close(); err != nil {
			r.logger.Error("关闭嵌入模型客户端失败", zap.Error(err))
		}
	}

	r.logger.Info("RAG系统资源已关闭")
	return nil
}

// QueryParams 查询参数
type QueryParams struct {
	Mode              string // "local", "global", "hybrid", "naive", "mix"
	OnlyNeedContext   bool
	OnlyNeedPrompt    bool
	ResponseType      string
	Stream            bool
	TopK              int
	MaxTokenForUnit   int
	MaxTokenForGlobal int
	MaxTokenForLocal  int
}
