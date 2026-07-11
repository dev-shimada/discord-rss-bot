package main

import (
	"errors"
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

	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// Discord Bot Token
	token := os.Getenv("DISCORD_BOT_TOKEN")

	db := database.NewDB()
	if db == nil {
		return errors.New("failed to initialize database")
	}
	defer database.CloseDB(db)

	// Create a new Discord session using the provided bot token.
	session, err := router.NewRouter(token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	// DI
	dh := di.DiscordHandler(db, session)

	// Open Discord session
	router.Open(session, dh)
	return nil
}
