package api

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:embed static/*
var staticFiles embed.FS

// setupWebRoutes 设置Web界面路由
func (s *Server) setupWebRoutes() {
	s.logger.Info("设置Web界面路由")

	// 提供静态文件服务
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		s.logger.Error("加载静态文件失败", zap.Error(err))
		return
	}

	// 设置静态文件服务
	s.router.StaticFS("/static", http.FS(staticFS))

	// 主页路由
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/static/index.html")
	})

	// Web界面API路由
	web := s.router.Group("/web")
	{
		// 获取配置信息
		web.GET("/config", s.handleGetConfig)

		// 更新配置信息
		web.POST("/config", s.handleUpdateConfig)
	}
}

// handleGetConfig 处理获取配置信息请求
func (s *Server) handleGetConfig(c *gin.Context) {
	// 返回部分配置信息（排除敏感信息如API密钥）
	c.JSON(http.StatusOK, gin.H{
		"rag": gin.H{
			"workingDir":     s.cfg.Rag.WorkingDir,
			"inputDir":       s.cfg.Rag.InputDir,
			"maxTokens":      s.cfg.Rag.MaxTokens,
			"maxEmbedTokens": s.cfg.Rag.MaxEmbedTokens,
			"chunkSize":      s.cfg.Rag.ChunkSize,
			"chunkOverlap":   s.cfg.Rag.ChunkOverlap,
			"topK":           s.cfg.Rag.TopK,
		},
		"llm": gin.H{
			"provider": s.cfg.Llm.Provider,
			"model":    s.cfg.Llm.Model,
			"baseUrl":  s.cfg.Llm.BaseUrl,
		},
		"embedding": gin.H{
			"provider": s.cfg.Embedding.Provider,
			"model":    s.cfg.Embedding.Model,
			"baseUrl":  s.cfg.Embedding.BaseUrl,
		},
		"storage": gin.H{
			"vectorType": s.cfg.Storage.VectorType,
			"kvType":     s.cfg.Storage.KvType,
			"graphType":  s.cfg.Storage.GraphType,
		},
	})
}

// handleUpdateConfig 处理更新配置信息请求
func (s *Server) handleUpdateConfig(c *gin.Context) {
	// 解析请求数据
	var configData struct {
		Rag struct {
			WorkingDir     string `json:"workingDir"`
			InputDir       string `json:"inputDir"`
			MaxTokens      int    `json:"maxTokens"`
			MaxEmbedTokens int    `json:"maxEmbedTokens"`
			ChunkSize      int    `json:"chunkSize"`
			ChunkOverlap   int    `json:"chunkOverlap"`
			TopK           int    `json:"topK"`
		} `json:"rag"`
		Llm struct {
			Provider string `json:"provider"`
			Model    string `json:"model"`
			ApiKey   string `json:"apiKey,omitempty"`
			BaseUrl  string `json:"baseUrl,omitempty"`
		} `json:"llm"`
		Embedding struct {
			Provider string `json:"provider"`
			Model    string `json:"model"`
			ApiKey   string `json:"apiKey,omitempty"`
			BaseUrl  string `json:"baseUrl,omitempty"`
		} `json:"embedding"`
		Storage struct {
			VectorType string `json:"vectorType"`
			KvType     string `json:"kvType"`
			GraphType  string `json:"graphType"`
		} `json:"storage"`
		Server struct {
			ApiKey string `json:"apiKey,omitempty"`
		} `json:"server,omitempty"`
	}

	if err := c.ShouldBindJSON(&configData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无效的配置数据: %v", err)})
		return
	}

	// 更新配置
	s.cfg.Rag.WorkingDir = configData.Rag.WorkingDir
	s.cfg.Rag.InputDir = configData.Rag.InputDir
	s.cfg.Rag.MaxTokens = configData.Rag.MaxTokens
	s.cfg.Rag.MaxEmbedTokens = configData.Rag.MaxEmbedTokens
	s.cfg.Rag.ChunkSize = configData.Rag.ChunkSize
	s.cfg.Rag.ChunkOverlap = configData.Rag.ChunkOverlap
	s.cfg.Rag.TopK = configData.Rag.TopK

	s.cfg.Llm.Provider = configData.Llm.Provider
	s.cfg.Llm.Model = configData.Llm.Model
	if configData.Llm.ApiKey != "" {
		s.cfg.Llm.ApiKey = configData.Llm.ApiKey
	}
	if configData.Llm.BaseUrl != "" {
		s.cfg.Llm.BaseUrl = configData.Llm.BaseUrl
	}

	s.cfg.Embedding.Provider = configData.Embedding.Provider
	s.cfg.Embedding.Model = configData.Embedding.Model
	if configData.Embedding.ApiKey != "" {
		s.cfg.Embedding.ApiKey = configData.Embedding.ApiKey
	}
	if configData.Embedding.BaseUrl != "" {
		s.cfg.Embedding.BaseUrl = configData.Embedding.BaseUrl
	}

	s.cfg.Storage.VectorType = configData.Storage.VectorType
	s.cfg.Storage.KvType = configData.Storage.KvType
	s.cfg.Storage.GraphType = configData.Storage.GraphType

	if configData.Server.ApiKey != "" {
		s.cfg.Server.ApiKey = configData.Server.ApiKey
	}

	// 保存配置到文件
	if err := saveConfigToFile(s.cfg); err != nil {
		s.logger.Error("保存配置失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("保存配置失败: %v", err)})
		return
	}

	// 重新初始化RAG系统
	s.initRAG()

	c.JSON(http.StatusOK, gin.H{"message": "配置已更新"})
}
