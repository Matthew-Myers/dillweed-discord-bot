package messagestore

const (
	OpenAIV1 = iota
	ChatML
)

type transcriptFormattor func(prompt string, msgs []MessageRecord) []OpenAIMessageFormat

type OpenAIMessageFormat struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// msgs should be ordered oldest -> most recent
var OpenAIV1Format transcriptFormattor = func(prompt string, msgs []MessageRecord) []OpenAIMessageFormat {
	openAITranscript := []OpenAIMessageFormat{
		{Role: "system", Content: prompt},
	}

	for _, msg := range msgs {
		openAITranscript = append(openAITranscript, OpenAIMessageFormat{Role: msg.AuthorName, Content: msg.Content})
	}
	return openAITranscript
}
