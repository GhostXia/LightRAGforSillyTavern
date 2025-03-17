package lightrag

import (
	"context"
)

// LLMClient 大语言模型客户端接口
type LLMClient interface {
	// Complete 完成文本生成
	Complete(ctx context.Context, prompt string, systemPrompt string, maxTokens int) (string, error)

	// CompleteStream 流式完成文本生成
	CompleteStream(ctx context.Context, prompt string, systemPrompt string, maxTokens int, callback func(string) error) error

	// Close 关闭客户端
	Close() error
}

// EmbeddingClient 嵌入模型客户端接口
type EmbeddingClient interface {
	// Embed 生成文本嵌入
	Embed(ctx context.Context, texts []string) ([][]float32, error)

	// Close 关闭客户端
	Close() error
}

// OpenAIMessage OpenAI消息格式
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIChatRequest OpenAI聊天请求
type OpenAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

// OpenAIEmbeddingRequest OpenAI嵌入请求
type OpenAIEmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// OpenAIChatResponse OpenAI聊天响应
type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// OpenAIEmbeddingResponse OpenAI嵌入响应
type OpenAIEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}
