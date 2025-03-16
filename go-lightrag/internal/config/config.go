package config

import (
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Rag       RagConfig       `mapstructure:"rag"`
	Llm       LlmConfig       `mapstructure:"llm"`
	Embedding EmbeddingConfig `mapstructure:"embedding"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Log       LogConfig       `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	ApiKey string `mapstructure:"api_key"`
}

// RagConfig RAG系统配置
type RagConfig struct {
	WorkingDir        string `mapstructure:"working_dir"`
	InputDir          string `mapstructure:"input_dir"`
	MaxTokens         int    `mapstructure:"max_tokens"`
	MaxEmbedTokens    int    `mapstructure:"max_embed_tokens"`
	ChunkSize         int    `mapstructure:"chunk_size"`
	ChunkOverlap      int    `mapstructure:"chunk_overlap"`
	TopK              int    `mapstructure:"top_k"`
	EntitySummarySize int    `mapstructure:"entity_summary_size"`
}

// LlmConfig LLM配置
type LlmConfig struct {
	Provider string `mapstructure:"provider"` // openai, ollama, etc.
	Model    string `mapstructure:"model"`
	ApiKey   string `mapstructure:"api_key"`
	BaseUrl  string `mapstructure:"base_url"`
}

// EmbeddingConfig 嵌入模型配置
type EmbeddingConfig struct {
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`
	ApiKey   string `mapstructure:"api_key"`
	BaseUrl  string `mapstructure:"base_url"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	VectorType string         `mapstructure:"vector_type"` // memory, postgres, milvus, etc.
	KvType     string         `mapstructure:"kv_type"`     // memory, postgres, mongo, etc.
	GraphType  string         `mapstructure:"graph_type"`  // memory, neo4j, postgres, etc.
	Neo4j      Neo4jConfig    `mapstructure:"neo4j"`
	Postgres   PostgresConfig `mapstructure:"postgres"`
	Mongo      MongoConfig    `mapstructure:"mongo"`
	Milvus     MilvusConfig   `mapstructure:"milvus"`
}

// Neo4jConfig Neo4j配置
type Neo4jConfig struct {
	Uri      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// PostgresConfig PostgreSQL配置
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SslMode  string `mapstructure:"ssl_mode"`
}

// MongoConfig MongoDB配置
type MongoConfig struct {
	Uri      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

// MilvusConfig Milvus配置
type MilvusConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 9621)

	// RAG默认配置
	viper.SetDefault("rag.working_dir", "./rag_storage")
	viper.SetDefault("rag.input_dir", "./inputs")
	viper.SetDefault("rag.max_tokens", 32768)
	viper.SetDefault("rag.max_embed_tokens", 8192)
	viper.SetDefault("rag.chunk_size", 1024)
	viper.SetDefault("rag.chunk_overlap", 128)
	viper.SetDefault("rag.top_k", 60)
	viper.SetDefault("rag.entity_summary_size", 4000)

	// LLM默认配置
	viper.SetDefault("llm.provider", "openai")
	viper.SetDefault("llm.model", "gpt-4o-mini")

	// 嵌入模型默认配置
	viper.SetDefault("embedding.provider", "openai")
	viper.SetDefault("embedding.model", "text-embedding-3-large")

	// 存储默认配置
	viper.SetDefault("storage.vector_type", "memory")
	viper.SetDefault("storage.kv_type", "memory")
	viper.SetDefault("storage.graph_type", "memory")

	// 日志默认配置
	viper.SetDefault("log.level", "info")
}
