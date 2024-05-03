package chat

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type Chat struct {
	Token   string
	Content string
	PreLoad string
	Model   string
}

type Response struct {
	Content string
	Error   string
}

func (r *Response) GetResponse() string {
	return r.Content
}

func (r *Response) GetError() string {
	return r.Error
}

func NewMessage(body *Chat) *Response {
	if body == nil {
		return &Response{
			Error: fmt.Sprintf("body can not be nil"),
		}
	}
	if body.Token == "" {
		return &Response{
			Error: fmt.Sprintf("token can not be nil"),
		}
	}
	var model = openai.GPT3Dot5Turbo
	if body.Model != "" {
		model = body.Model
	}
	client := openai.NewClient(body.Token)
	if body.PreLoad != "" {
		body.Content = fmt.Sprintf("the user input is by the user and the proloaded info is to help chat gpt understand what the user wants. preloaded info: %s user input: %s", body.PreLoad, body.Content)
	}
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: body.Content,
				},
			},
		},
	)
	if err != nil {
		return &Response{
			Error: fmt.Sprintf("ChatCompletion error: %v\n", err),
		}
	}

	return &Response{
		Content: resp.Choices[0].Message.Content,
	}
}
