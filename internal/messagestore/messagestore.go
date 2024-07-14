package messagestore

import (
	"internal/botconfig"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

// Define a struct to represent your table's schema
type MessageRecord struct {
	Id          string                         `json:"Id"`
	CreatedAt   string                         `json:"CreatedAt"`
	AuthorId    string                         `json:"AuthorId"`
	GlobalName  string                         `json:"GlobalName"`
	Content     string                         `json:"Content"`
	Attachments []*discordgo.MessageAttachment `json:"Attachments"`
	AuthorName  string                         `json:"AuthorName"`
	// Add other fields as needed
}

type MessageStore struct {
	messageCache *MessageCache
	config       *botconfig.AppConfig
}

func NewMessageStore(cacheLimit int, config *botconfig.AppConfig) *MessageStore {
	return &MessageStore{
		messageCache: NewMessageCache(cacheLimit, config),
		config:       config,
	}
}

func (ms *MessageStore) store(s *discordgo.Session, m *discordgo.MessageCreate) {

	var AuthorName = m.Author.Username
	if m.Author.ID == s.State.User.ID {
		AuthorName = "Dillweed"
	}
	msgRec := MessageRecord{
		Id:          m.ID,
		CreatedAt:   m.Timestamp.Format("2006-01-02T15:04:05Z"),
		AuthorId:    m.Author.ID,
		GlobalName:  m.Author.GlobalName,
		AuthorName:  AuthorName,
		Content:     m.Content,
		Attachments: m.Attachments,
	}
	// Strore the messsage in postgres
	if ms.config.SQLSettings.UseSQL {
		println("SQL is not yet supproted")
	}
	// Update in memory cache for the channel message log
	//   To handle more than just a few handfuls of servers + channels, this will need some work
	//   Given this is going to be for a single server, this will be fine.
	//   Max char count of a discord message is 2000 (ignoring nitro)
	//   Worst case, 120 messages will end up consuming about 1GB of ram.
	ms.messageCache.Put(m.GuildID, m.ChannelID, msgRec)
}

func (ms *MessageStore) MessageCreateEventHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// If the message is a command, ignore

	// store message and update the cache
	ms.store(s, m)

	if ms.respond(s, m) {
		channelStore, hit := ms.messageCache.ChannelStore(m.GuildID, m.ChannelID)
		if !hit {
			// TODO: actually handle
			log.Fatal("No messages in the cache,")
		}
		promptIndex, prompt := ms.config.Prompt.ChoosePrompt(channelStore.personalityIndex)
		channelStore.personalityIndex = promptIndex
		replyString := LLMRequest(prompt, channelStore.messages, ms.config)
		_, err := s.ChannelMessageSend(m.ChannelID, replyString)
		if err != nil {
			println(err.Error())
		}
	}
}

func (ms *MessageStore) respond(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Author.ID == s.State.User.ID {
		return false
	}

	return rand.Intn(int(ms.config.Prompt.Settings.Roll)) == 0
}
