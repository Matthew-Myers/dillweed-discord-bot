package messagestore

import (
	"bytes"
	"encoding/json"
	"internal/botconfig"
	"log"
	"net/http"
	"strings"
)

type OpenAIRequest struct {
	botconfig.LLMMessageConfig
	Messages []OpenAIMessageFormat `json:"messages"`
}
type OpenAIResponseMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}
type OpenAIChoice struct {
	FinishReason string                `json:"finish_reason"`
	Index        int                   `json:"index"`
	Message      OpenAIResponseMessage `json:"message"`
}
type OpenAIUsage struct {
	CompletionTokeens int `json:"completion_tokens"`
	PromptTokens      int `json:"prompt_tokens"`
	TotalTokens       int `json:"total_tokens"`
}
type OpenAIResponseBody struct {
	Choices []OpenAIChoice `json:"choices"`
	Created int32          `json:"created"`
	Model   string         `json:"model"`
	Object  string         `json:"chat.completion"`
	Usage   OpenAIUsage    `json:"usage"`
	Id      string         `json:"id"`
}

func LLMRequest(prompt string, msgs []MessageRecord, config *botconfig.AppConfig) string {

	transcript := OpenAIV1Format(prompt, msgs)

	llmConfig := config.LLMMessageConfig
	// TODO: move this config to an on read
	requestStruct := OpenAIRequest{
		llmConfig,
		transcript,
	}

	postBody, _ := json.Marshal(requestStruct)

	reqBytes := bytes.NewBuffer(postBody)
	// TODO: move the url to a config
	resp, err := http.Post(config.LLMHost.URL, "application/json", reqBytes)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	var oresp OpenAIResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&oresp); err != nil {
		log.Fatal("Could not decode the response body")
	}
	// llama cpp keep adding the instruction end at the end of the message and dillweed in the prefix.  Just going to cut it here
	return formatLLMResponse(oresp.Choices[0].Message.Content)
}

func formatLLMResponse(message string) string {
	var msg = message
	prefix := "Dillweed: "
	if strings.HasPrefix(message, prefix) {
		msg = strings.TrimPrefix(message, prefix)
	}
	prefix = "Dillweed:\n"
	if strings.HasPrefix(message, prefix) {
		msg = strings.TrimPrefix(message, prefix)
	}
	msg = strings.Replace(msg, "<|im_end|>", "", -1)
	return msg
}
