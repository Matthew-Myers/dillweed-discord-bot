module dillweed

go 1.22.4

require internal/messagestore v1.0.0
replace internal/messagestore => ./internal/messagestore

require (
	github.com/bwmarrin/discordgo v0.28.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)
