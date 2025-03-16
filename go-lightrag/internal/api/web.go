package api

import (
	"embed"
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
		},
		"embedding": gin.H{
			"provider": s.cfg.Embedding.Provider,
			"model":    s.cfg.Embedding.Model,
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
	// TODO: 实现配置更新逻辑
	c.JSON(http.StatusOK, gin.H{"message": "配置更新功能待实现"})
}
