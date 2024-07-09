package messagestore

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type OpenAIRequest struct {
	Continue                 bool                  `json:"_continue"`
	Add_bos_token            bool                  `json:"add_bos_token"`
	Ban_eos_token            bool                  `json:"ban_eos_token"`
	Max_tokens               int                   `json:"max_tokens"`
	Custom_token_bans        string                `json:"custom_token_bans"`
	Do_sample                bool                  `json:"do_sample"`
	Early_stopping           bool                  `json:"early_stopping"`
	Epsilon_cutoff           float32               `json:"epsilon_cutoff"`
	Eta_cutoff               float32               `json:"eta_cutoff"`
	Grammar_string           string                `json:"grammar_string"`
	Guidance_scale           float32               `json:"guidance_scale"`
	Length_penalty           float32               `json:"length_penalty"`
	Min_length               float32               `json:"min_length"`
	Mirostat_eta             float32               `json:"mirostat_eta"`
	Mirostat_mode            float32               `json:"mirostat_mode"`
	Mirostat_tau             int                   `json:"mirostat_tau"`
	Negative_prompt          string                `json:"negative_prompt"`
	No_repeat_ngram_size     float32               `json:"no_repeat_ngram_size"`
	Num_beams                int                   `json:"num_beams"`
	Penalty_alpha            float32               `json:"penalty_alpha"`
	Preset                   string                `json:"preset"`
	Regenerate               bool                  `json:"regenerate"`
	Repetition_penalty       float32               `json:"repetition_penalty"`
	Repetition_penalty_range float32               `json:"repetition_penalty_range"`
	Seed                     int                   `json:"seed"`
	Skip_special_tokens      bool                  `json:"skip_special_tokens"`
	Stopping_strings         []string              `json:"stopping_strings"`
	Temperature              float32               `json:"temperature"`
	Tfs                      float32               `json:"tfs"`
	Top_a                    float32               `json:"top_a"`
	Top_k                    int                   `json:"top_k"`
	Top_p                    float32               `json:"top_p"`
	Typical_p                float32               `json:"typical_p"`
	Instruction_template     string                `json:"instruction_template"`
	Messages                 []OpenAIMessageFormat `json:"messages"`
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

func LLMRequest(prompt string, msgs []MessageRecord) string {
	transcript := OpenAIV1Format(prompt, msgs)

	// TODO: move this config to an on read
	var requestStruct OpenAIRequest = OpenAIRequest{
		Continue:                 false,
		Add_bos_token:            true,
		Ban_eos_token:            false,
		Max_tokens:               1024,
		Custom_token_bans:        "",
		Do_sample:                true,
		Early_stopping:           false,
		Epsilon_cutoff:           0,
		Eta_cutoff:               0,
		Grammar_string:           "",
		Guidance_scale:           1,
		Length_penalty:           1,
		Min_length:               0,
		Mirostat_eta:             0.1,
		Mirostat_mode:            0,
		Mirostat_tau:             5,
		Negative_prompt:          "",
		No_repeat_ngram_size:     0,
		Num_beams:                1,
		Penalty_alpha:            0,
		Regenerate:               false,
		Repetition_penalty:       1.18,
		Repetition_penalty_range: 0,
		Seed:                     -1,
		Skip_special_tokens:      true,
		Temperature:              0.8,
		Tfs:                      1,
		Top_a:                    0,
		Top_k:                    50,
		Top_p:                    1,
		Typical_p:                1,
		Messages:                 transcript,
	}

	postBody, _ := json.Marshal(requestStruct)
	reqBytes := bytes.NewBuffer(postBody)
	println(postBody)
	// TODO: move the url to a config
	resp, err := http.Post("http://192.168.1.153:9000/v1/chat/completions", "application/json", reqBytes)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	var oresp OpenAIResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&oresp); err != nil {
		log.Fatal("Could not decode the response body")
	}
	// llama cpp keep adding the instruction end at the end of the message.  Just going to cut it here
	return strings.Replace(oresp.Choices[0].Message.Content, "<|im_end|>", "", -1)
}
