package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HerSophia/go-lightrag/internal/config"
	"github.com/HerSophia/go-lightrag/pkg/lightrag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Server 表示API服务器
type Server struct {
	cfg    *config.Config
	logger *zap.Logger
	router *gin.Engine
	server *http.Server
	rag    *lightrag.LightRAG
}

// NewServer 创建新的API服务器
func NewServer(cfg *config.Config, logger *zap.Logger) *Server {
	// 设置Gin模式
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// 创建服务器
	s := &Server{
		cfg:    cfg,
		logger: logger,
		router: router,
	}

	// 初始化RAG
	s.initRAG()

	// 设置路由
	s.setupRoutes()

	return s
}

// initRAG 初始化RAG系统
func (s *Server) initRAG() {
	// 创建LightRAG实例
	s.rag = lightrag.NewLightRAG(
		s.cfg.Rag.WorkingDir,
		s.cfg.Rag.MaxTokens,
		s.cfg.Rag.MaxEmbedTokens,
		s.cfg.Rag.ChunkSize,
		s.cfg.Rag.ChunkOverlap,
		s.cfg.Rag.TopK,
		s.logger,
	)

	// 初始化存储
	s.rag.InitStorage(s.cfg.Storage)

	// 初始化LLM和嵌入模型
	s.rag.InitLLM(s.cfg.Llm)
	s.rag.InitEmbedding(s.cfg.Embedding)
}

// setupRoutes 设置API路由
func (s *Server) setupRoutes() {
	// API密钥中间件
	authMiddleware := func(c *gin.Context) {
		if s.cfg.Server.ApiKey != "" {
			apiKey := c.GetHeader("Authorization")
			if apiKey != fmt.Sprintf("Bearer %s", s.cfg.Server.ApiKey) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "无效的API密钥",
				})
				return
			}
		}
		c.Next()
	}

	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API路由组
	api := s.router.Group("/api")
	api.Use(authMiddleware)
	{
		// 查询接口
		api.POST("/query", s.handleQuery)

		// 插入文本接口
		api.POST("/insert", s.handleInsert)

		// 上传文件接口
		api.POST("/upload", s.handleUpload)

		// 获取文档状态接口
		api.GET("/documents", s.handleGetDocuments)

		// 删除文档接口
		api.DELETE("/documents/:id", s.handleDeleteDocument)

		// 获取知识图谱列表
		api.GET("/graphs", s.handleGetGraphs)

		// 获取知识图谱数据
		api.GET("/graphs/:id", s.handleGetGraph)
	}

	// OpenAI兼容API路由组
	openai := s.router.Group("/v1")
	openai.Use(authMiddleware)
	{
		// 兼容OpenAI的聊天完成接口
		openai.POST("/chat/completions", s.handleChatCompletions)

		// 兼容OpenAI的嵌入接口
		openai.POST("/embeddings", s.handleEmbeddings)
	}

	// 设置Web界面路由
	s.setupWebRoutes()
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.server.Shutdown(ctx)
	}

	// 关闭RAG系统
	if s.rag != nil {
		s.rag.Close()
	}
}

// 处理函数 - 这些将在后续实现
func (s *Server) handleQuery(c *gin.Context) {
	// TODO: 实现查询处理
	c.JSON(http.StatusOK, gin.H{"message": "查询功能待实现"})
}

func (s *Server) handleInsert(c *gin.Context) {
	// TODO: 实现插入处理
	c.JSON(http.StatusOK, gin.H{"message": "插入功能待实现"})
}

func (s *Server) handleUpload(c *gin.Context) {
	// TODO: 实现文件上传处理
	c.JSON(http.StatusOK, gin.H{"message": "上传功能待实现"})
}

func (s *Server) handleGetDocuments(c *gin.Context) {
	// TODO: 实现获取文档列表
	c.JSON(http.StatusOK, gin.H{"message": "获取文档列表功能待实现"})
}

func (s *Server) handleDeleteDocument(c *gin.Context) {
	// TODO: 实现删除文档
	c.JSON(http.StatusOK, gin.H{"message": "删除文档功能待实现"})
}

// handleGetGraphs 处理获取知识图谱列表请求
func (s *Server) handleGetGraphs(c *gin.Context) {
	// 从工作目录中获取所有知识图谱
	graphs := []gin.H{}

	// TODO: 实现从存储中获取图谱列表
	// 示例数据
	graphs = append(graphs, gin.H{
		"id":   "example1",
		"name": "示例知识图谱1",
	})

	c.JSON(http.StatusOK, gin.H{
		"graphs": graphs,
	})
}

// handleGetGraph 处理获取知识图谱数据请求
func (s *Server) handleGetGraph(c *gin.Context) {
	graphID := c.Param("id")

	// TODO: 实现从存储中获取图谱数据
	// 示例数据
	nodes := []gin.H{
		{
			"id":    "1",
			"label": "节点1",
			"color": "#ff0000",
		},
		{
			"id":    "2",
			"label": "节点2",
			"color": "#00ff00",
		},
		{
			"id":    "3",
			"label": "节点3",
			"color": "#0000ff",
		},
	}

	edges := []gin.H{
		{
			"from":  "1",
			"to":    "2",
			"label": "关系1",
		},
		{
			"from":  "2",
			"to":    "3",
			"label": "关系2",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    graphID,
		"name":  "示例知识图谱" + graphID,
		"nodes": nodes,
		"edges": edges,
	})
}

// handleChatCompletions 处理OpenAI兼容的聊天完成请求
func (s *Server) handleChatCompletions(c *gin.Context) {
	var req lightrag.OpenAIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无效的请求: %v", err)})
		return
	}

	// 验证请求
	if req.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型名称不能为空"})
		return
	}

	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息不能为空"})
		return
	}

	// 提取用户消息和系统消息
	var userMessage, systemMessage string
	for _, msg := range req.Messages {
		if msg.Role == "user" {
			userMessage = msg.Content
		} else if msg.Role == "system" {
			systemMessage = msg.Content
		}
	}

	if userMessage == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户消息不能为空"})
		return
	}

	// 设置最大token数
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = s.cfg.Rag.MaxTokens
	}

	// 处理流式请求
	if req.Stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Transfer-Encoding", "chunked")

		clientGone := c.Writer.CloseNotify()
		c.Stream(func(w io.Writer) bool {
			// 使用RAG系统处理请求
			err := s.rag.CompleteStream(c.Request.Context(), userMessage, systemMessage, maxTokens, func(chunk string) error {
				// 检查客户端是否已断开连接
				select {
				case <-clientGone:
					return errors.New("客户端已断开连接")
				default:
					// 继续处理
				}

				// 构造SSE消息
				sseMsg := fmt.Sprintf("data: {\"id\":\"%s\",\"object\":\"chat.completion.chunk\",\"created\":%d,\"model\":\"%s\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"%s\"},\"finish_reason\":null}]}\n\n",
					"chatcmpl-"+uuid.New().String(), time.Now().Unix(), req.Model, strings.ReplaceAll(chunk, "\"", "\\\""))

				// 发送SSE消息
				_, err := fmt.Fprint(w, sseMsg)
				return err
			})

			if err != nil {
				s.logger.Error("流式处理请求失败", zap.Error(err))
			}

			// 发送结束消息
			endMsg := "data: {\"id\":\"chatcmpl-" + uuid.New().String() + "\",\"object\":\"chat.completion.chunk\",\"created\":" +
				fmt.Sprintf("%d", time.Now().Unix()) + ",\"model\":\"" + req.Model + "\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n"
			fmt.Fprint(w, endMsg)
			fmt.Fprint(w, "data: [DONE]\n\n")

			return false
		})
		return
	}

	// 非流式请求
	response, err := s.rag.Query(userMessage, lightrag.QueryParams{
		MaxTokenForGlobal: maxTokens,
	})

	if err != nil {
		s.logger.Error("处理请求失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("处理请求失败: %v", err)})
		return
	}

	// 构造OpenAI兼容的响应
	openaiResp := lightrag.OpenAIChatResponse{
		ID:      "chatcmpl-" + uuid.New().String(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []struct {
			Index        int                    `json:"index"`
			Message      lightrag.OpenAIMessage `json:"message"`
			FinishReason string                 `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: lightrag.OpenAIMessage{
					Role:    "assistant",
					Content: response,
				},
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     len(userMessage) / 4, // 简单估算
			CompletionTokens: len(response) / 4,    // 简单估算
			TotalTokens:      (len(userMessage) + len(response)) / 4,
		},
	}

	c.JSON(http.StatusOK, openaiResp)
}

// handleEmbeddings 处理OpenAI兼容的嵌入请求
func (s *Server) handleEmbeddings(c *gin.Context) {
	var req lightrag.OpenAIEmbeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无效的请求: %v", err)})
		return
	}

	// 验证请求
	if req.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型名称不能为空"})
		return
	}

	if len(req.Input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "输入文本不能为空"})
		return
	}

	// 生成嵌入
	embeddings, err := s.rag.Embed(c.Request.Context(), req.Input)
	if err != nil {
		s.logger.Error("生成嵌入失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("生成嵌入失败: %v", err)})
		return
	}

	// 构造OpenAI兼容的响应
	openaiResp := lightrag.OpenAIEmbeddingResponse{
		Object: "list",
		Data: make([]struct {
			Object    string    `json:"object"`
			Embedding []float32 `json:"embedding"`
			Index     int       `json:"index"`
		}, len(embeddings)),
		Model: req.Model,
		Usage: struct {
			PromptTokens int `json:"prompt_tokens"`
			TotalTokens  int `json:"total_tokens"`
		}{
			PromptTokens: 0,
			TotalTokens:  0,
		},
	}

	// 计算token使用量（简单估算）
	for _, text := range req.Input {
		openaiResp.Usage.PromptTokens += len(text) / 4
	}
	openaiResp.Usage.TotalTokens = openaiResp.Usage.PromptTokens

	// 填充嵌入数据
	for i, embedding := range embeddings {
		openaiResp.Data[i] = struct {
			Object    string    `json:"object"`
			Embedding []float32 `json:"embedding"`
			Index     int       `json:"index"`
		}{
			Object:    "embedding",
			Embedding: embedding,
			Index:     i,
		}
	}

	c.JSON(http.StatusOK, openaiResp)
}
