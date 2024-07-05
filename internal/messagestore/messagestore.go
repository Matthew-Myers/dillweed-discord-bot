package messagestore

// Define a struct to represent your table's schema
type MessageRecord struct {
	ID          int
	CreatedAt   string
	AuthorId    string
	GlobalName  string
	Content     string
	Attachments []string
	AuthorName  string
	// Add other fields as needed
}
