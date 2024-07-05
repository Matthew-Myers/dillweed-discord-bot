package messagestore

import "sort"

const (
	OpenAIV1 = iota
	ChatML
)

type transcriptFormattor func(prompt string, msgs []MessageRecord) any

type openAIMessageFormat struct {
	role    string
	content string
}

var openAIV1Format transcriptFormattor = func(prompt string, msgs []MessageRecord) any {
	openAITranscript := []openAIMessageFormat{
		{role: "system", content: prompt},
	}

	// Should be in order, but just in case
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].CreatedAt < msgs[j].CreatedAt
	})

	for _, msg := range msgs {
		openAITranscript = append(openAITranscript, openAIMessageFormat{role: msg.AuthorName, content: msg.Content})
	}

	return openAITranscript
}
