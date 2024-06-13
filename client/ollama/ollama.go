package ollama

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dmarkham/agency/client"
	"github.com/ollama/ollama/api"
)

type Model string

type Client struct {
	client *api.Client
}

func New(baseURL string, c *http.Client) *Client {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}
	client := api.NewClient(u, c)

	return &Client{
		client: client,
	}
}

func NewClient(cl *api.Client) *Client {
	return &Client{
		client: cl,
	}
}

func (Client) SupportsStreaming() bool {
	return true
}

var _ client.Client = (*Client)(nil)

func maybeWrapError(err error) error {
	//TODO Double check this to see if this is how it's intended to be used
	e := &api.StatusError{}
	if errors.As(err, &e) && e.StatusCode == 429 || e.StatusCode == 500 {
		return client.Retryable(err)
	}
	return err

}

func (cl *Client) CreateChatCompletion(ctx context.Context, request client.ChatCompletionRequest) (client.ChatCompletionResponse, error) {
	req, err := TranslateRequest(request)
	if err != nil {
		return client.ChatCompletionResponse{}, err
	}
	if request.Stream == nil {
		req.Stream = new(bool)
	}

	resp := &api.ChatResponse{}
	finalMessage := &strings.Builder{}
	if request.Stream != nil {
		request.Stream = io.MultiWriter(request.Stream, finalMessage)
	} else {
		request.Stream = io.MultiWriter(finalMessage)
	}

	respFunc := func(r api.ChatResponse) error {
		resp = &r
		fmt.Fprint(request.Stream, r.Message.Content)
		return nil
	}

	err = cl.client.Chat(ctx, req, respFunc)
	if err != nil {
		return client.ChatCompletionResponse{}, maybeWrapError(err)
	}
	resp.Message.Content = finalMessage.String()
	//fmt.Print("HERE", resp)

	return TranslateResponse(*resp), nil
}

var roleMapping = map[client.Role]string{
	client.User:      "user",
	client.System:    "system",
	client.Assistant: "assistant",
}

func TranslateRequest(clientReq client.ChatCompletionRequest) (*api.ChatRequest, error) {

	req := &api.ChatRequest{
		Model: clientReq.Model,
	}
	if clientReq.Stream == nil {
		req.Stream = new(bool)
	}

	req.Options = clientReq.CustomParams

	for _, message := range clientReq.Messages {
		req.Messages = append(req.Messages, api.Message{
			Content: message.Content,
			Role:    roleMapping[message.Role],
		})
	}

	return req, nil

}

// TranslateResponse translates a ChatCompletionResponse from the openai package to one from the client package
func TranslateResponse(apiResp api.ChatResponse) client.ChatCompletionResponse {
	// Create a new slice to hold the translated messages
	clientMessages := make([]client.Message, 0)

	// Translate each choice's message to the client's Message type
	clientMessages = append(clientMessages, client.Message{
		Content: apiResp.Message.Content,
		Role:    client.Role(apiResp.Message.Role),
	})

	// Return a new ChatCompletionResponse from the client package, using the translated messages
	return client.ChatCompletionResponse{
		Choices: clientMessages,
	}
}
