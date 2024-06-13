package ollama

import (
	"testing"

	"github.com/ollama/ollama/api"
	"github.com/dmarkham/agency/client"
	"github.com/stretchr/testify/assert"
)

func TestTranslateResponse(t *testing.T) {
	apiResponse := api.ChatResponse{
		Model: "test-model",
		Message: api.Message{
			Role:    "assistant",
			Content: "Hello, how can I assist you today?",
		},
	}

	expectedClientResponse := client.ChatCompletionResponse{
		Choices: []client.Message{
			{
				Content: "Hello, how can I assist you today?",
				Role:    "assistant",
			},
		},
	}

	actualClientResponse := TranslateResponse(apiResponse)

	if len(actualClientResponse.Choices) != len(expectedClientResponse.Choices) {
		t.Errorf("Expected number of choices %v, but got %v", len(expectedClientResponse.Choices), len(actualClientResponse.Choices))
	}

	for i := range actualClientResponse.Choices {
		if actualClientResponse.Choices[i].Content != expectedClientResponse.Choices[i].Content {
			t.Errorf("Expected content '%s', but got '%s'", expectedClientResponse.Choices[i].Content, actualClientResponse.Choices[i].Content)
		}
		if actualClientResponse.Choices[i].Role != expectedClientResponse.Choices[i].Role {
			t.Errorf("Expected role '%s', but got '%s'", expectedClientResponse.Choices[i].Role, actualClientResponse.Choices[i].Role)
		}
	}
}

func TestTranslateRequest(t *testing.T) {
	// Setup
	clientReq := client.ChatCompletionRequest{
		Model:       "gpt-4",
		MaxTokens:   100,
		Temperature: 0.8,
		Messages: []client.Message{
			{
				Content: "you are an assistant",
				Role:    client.System,
			},
			{
				Content: "How much is 2 + 2?",
				Role:    client.User,
			},
		},
		CustomParams: map[string]interface{}{
			"top_p":             float32(0.9),
			"stop":              []string{"\n"},
			"presence_penalty":  float32(0.6),
			"frequency_penalty": float32(0.5),
			"logit_bias":        map[string]int{"50256": -100},
			"user":              "user123",
		},
	}

	// Expected result
	expected := api.ChatRequest{
		Model: "gpt-4",
		Messages: []api.Message{
			{
				Content: "you are an assistant",
				Role:    "system",
			},
			{
				Content: "How much is 2 + 2?",
				Role:    "user",
			},
		},
	}

	// Run the function
	res, err := TranslateRequest(clientReq)

	// Validate results
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}
