package messagestore

import "github.com/bwmarrin/discordgo"

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

const isSQLConfigured = false

func StoreMessage(m *discordgo.MessageCreate, cache *MessageCache) {

	msgRec := MessageRecord{
		Id:          m.ID,
		CreatedAt:   m.Timestamp.Format("2006-01-02T15:04:05Z"),
		AuthorId:    m.Author.ID,
		GlobalName:  m.Author.GlobalName,
		AuthorName:  m.Author.Username,
		Content:     m.Content,
		Attachments: m.Attachments,
	}
	// Strore the messsage in postgres
	if isSQLConfigured {
		// store in SQL
	}

	// Update in memory cache for the channel message log
	//   To handle more than just a few handfuls of servers + channels, this will need some work
	//   Given this is going to be for a single server, this will be fine.
	//   Max char count of a discord message is 2000 (ignoring nitro)
	//   Worst case, 120 messages will end up consuming about 1GB of ram.
	cache.Put(m.GuildID, m.ChannelID, msgRec)

}
