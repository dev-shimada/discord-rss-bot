package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dev-shimada/discord-rss-bot/di"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/database"
	"github.com/dev-shimada/discord-rss-bot/router"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Discord Bot Token
	token := os.Getenv("DISCORD_BOT_TOKEN")

	db := database.NewDB()
	defer database.CloseDB(db)

	// Create a new Discord session using the provided bot token.
	session, err := router.NewRouter(token)
	if err != nil {
		slog.Error(fmt.Sprintf("error creating Discord session: %v", err))
	}

	// DI
	dh := di.DiscordHandler(db, session)

	// Open Discord session
	router.Open(session, dh)
}
