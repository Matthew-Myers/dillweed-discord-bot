package messagestore

import (
	"reflect"
	"testing"
)

var RECS = []MessageRecord{
	{
		ID:          0,
		CreatedAt:   "2022-07-05T12:24:36Z",
		AuthorId:    "uuid",
		GlobalName:  "EarliestMessage",
		Content:     "EMContents",
		Attachments: []string{"someAttachment"},
		AuthorName:  "AuthorEarliest",
	},
	{
		ID:          1,
		CreatedAt:   "2024-07-05T12:24:36Z",
		AuthorId:    "uuid",
		GlobalName:  "LatestMessage",
		Content:     "LMConents",
		Attachments: []string{"someAttachment"},
		AuthorName:  "AuthorLatest",
	},
	{
		ID:          2,
		CreatedAt:   "2023-07-05T12:24:36Z",
		AuthorId:    "uuid",
		GlobalName:  "MiddleMessage",
		Content:     "MMContents",
		Attachments: []string{"someAttachment"},
		AuthorName:  "AuthorMiddle",
	},
}

func TestOpenAiSystemPrompt(t *testing.T) {
	sysPrompt := "system prompt"
	var testArr []MessageRecord
	var response = openAIV1Format(sysPrompt, testArr)
	expected := []openAIMessageFormat{
		{role: "system", content: sysPrompt},
	}

	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Unexpected result. Expected:\n%v\nGot:\n%v", expected, response)
	}
}

func TestOpenAIFormatHandlesDateOrdering(t *testing.T) {
	sysPrompt := "system prompt"
	testArr := make([]MessageRecord, len(RECS))
	copy(testArr, RECS)
	var response = openAIV1Format(sysPrompt, testArr)
	expected := []openAIMessageFormat{
		{role: "system", content: sysPrompt},
		{role: RECS[0].AuthorName, content: RECS[0].Content},
		{role: RECS[2].AuthorName, content: RECS[2].Content},
		{role: RECS[1].AuthorName, content: RECS[1].Content},
	}

	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Unexpected result. Expected:\n%v\nGot:\n%v", expected, response)
	}
}
