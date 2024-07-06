package messagestore

const (
	OpenAIV1 = iota
	ChatML
)

type transcriptFormattor func(prompt string, msgs []MessageRecord) any

type openAIMessageFormat struct {
	role    string
	content string
}

// msgs should be ordered oldest -> most recent
var openAIV1Format transcriptFormattor = func(prompt string, msgs []MessageRecord) any {
	openAITranscript := []openAIMessageFormat{
		{role: "system", content: prompt},
	}

	for _, msg := range msgs {
		openAITranscript = append(openAITranscript, openAIMessageFormat{role: msg.AuthorName, content: msg.Content})
	}

	return openAITranscript
}
