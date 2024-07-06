package messagestore

import "github.com/bwmarrin/discordgo"

// Define a struct to represent your table's schema
type MessageRecord struct {
	Id          string
	CreatedAt   string
	AuthorId    string
	GlobalName  string
	Content     string
	Attachments []*discordgo.MessageAttachment
	AuthorName  string
	// Add other fields as needed
}

func StoreMessage(msg MessageRecord) {

}
