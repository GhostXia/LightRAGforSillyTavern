package embedding

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient 实现EmbeddingClient接口的OpenAI客户端
type OpenAIClient struct {
	client *openai.Client
	model  string
}

// NewOpenAIClient 创建新的OpenAI嵌入客户端
func NewOpenAIClient(apiKey, baseURL, model string) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, errors.New("OpenAI API密钥不能为空")
	}

	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	return &OpenAIClient{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}, nil
}

// Embed 生成文本嵌入
func (c *OpenAIClient) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, errors.New("嵌入文本不能为空")
	}

	req := openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(c.model),
		Input: texts,
	}

	resp, err := c.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("OpenAI返回的嵌入数据为空")
	}

	// 提取嵌入向量
	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}

// Close 关闭客户端
func (c *OpenAIClient) Close() error {
	// OpenAI客户端不需要特殊关闭操作
	return nil
}
