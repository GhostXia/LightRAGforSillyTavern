package llm

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient 实现LLMClient接口的OpenAI客户端
type OpenAIClient struct {
	client *openai.Client
	model  string
}

// NewOpenAIClient 创建新的OpenAI客户端
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

// Complete 完成文本生成
func (c *OpenAIClient) Complete(ctx context.Context, prompt string, systemPrompt string, maxTokens int) (string, error) {
	messages := []openai.ChatCompletionMessage{}

	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	req := openai.ChatCompletionRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: maxTokens,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("OpenAI返回的选择为空")
	}

	return resp.Choices[0].Message.Content, nil
}

// CompleteStream 流式完成文本生成
func (c *OpenAIClient) CompleteStream(ctx context.Context, prompt string, systemPrompt string, maxTokens int, callback func(string) error) error {
	messages := []openai.ChatCompletionMessage{}

	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	req := openai.ChatCompletionRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: maxTokens,
		Stream:    true,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		resp, err := stream.Recv()
		if errors.Is(err, openai.ErrStreamClosed) {
			break
		}
		if err != nil {
			return err
		}

		if len(resp.Choices) > 0 {
			content := resp.Choices[0].Delta.Content
			if content != "" {
				if err := callback(content); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Close 关闭客户端
func (c *OpenAIClient) Close() error {
	// OpenAI客户端不需要特殊关闭操作
	return nil
}
