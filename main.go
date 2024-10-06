package main

import (
	"log"
	"os"

	"github.com/dev-shimada/discord-rss-bot/di"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/database"
	"github.com/dev-shimada/discord-rss-bot/router"
)

func main() {
	// Discord Bot Token
	token := os.Getenv("DISCORD_BOT_TOKEN")

	db := database.NewDB()
	defer database.CloseDB(db)

	// Create a new Discord session using the provided bot token.
	session, err := router.NewRouter(token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	// DI
	dh := di.DiscordHandler(db, session)

	// Open Discord session
	router.Open(session, dh)
}
