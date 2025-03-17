package api

import (
	"github.com/HerSophia/go-lightrag/internal/config"
	"github.com/spf13/viper"
)

// saveConfigToFile 保存配置到文件
func saveConfigToFile(cfg *config.Config) error {
	// 设置配置值
	viper.Set("server.host", cfg.Server.Host)
	viper.Set("server.port", cfg.Server.Port)
	viper.Set("server.api_key", cfg.Server.ApiKey)

	viper.Set("rag.working_dir", cfg.Rag.WorkingDir)
	viper.Set("rag.input_dir", cfg.Rag.InputDir)
	viper.Set("rag.max_tokens", cfg.Rag.MaxTokens)
	viper.Set("rag.max_embed_tokens", cfg.Rag.MaxEmbedTokens)
	viper.Set("rag.chunk_size", cfg.Rag.ChunkSize)
	viper.Set("rag.chunk_overlap", cfg.Rag.ChunkOverlap)
	viper.Set("rag.top_k", cfg.Rag.TopK)
	viper.Set("rag.entity_summary_size", cfg.Rag.EntitySummarySize)

	viper.Set("llm.provider", cfg.Llm.Provider)
	viper.Set("llm.model", cfg.Llm.Model)
	viper.Set("llm.api_key", cfg.Llm.ApiKey)
	viper.Set("llm.base_url", cfg.Llm.BaseUrl)

	viper.Set("embedding.provider", cfg.Embedding.Provider)
	viper.Set("embedding.model", cfg.Embedding.Model)
	viper.Set("embedding.api_key", cfg.Embedding.ApiKey)
	viper.Set("embedding.base_url", cfg.Embedding.BaseUrl)

	viper.Set("storage.vector_type", cfg.Storage.VectorType)
	viper.Set("storage.kv_type", cfg.Storage.KvType)
	viper.Set("storage.graph_type", cfg.Storage.GraphType)

	// 保存到文件
	return viper.WriteConfig()
}
