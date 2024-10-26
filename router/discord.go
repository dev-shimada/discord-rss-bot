package router

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-shimada/discord-rss-bot/interface/discord"
)

func NewRouter(token string) (*discordgo.Session, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return dg, nil
}

func Open(dg *discordgo.Session, dh discord.DiscordHandler) {
	err := dg.Open()
	if err != nil {
		slog.Error(fmt.Sprintf("error opening connection: %v", err))
		return
	}
	defer dg.Close()

	// add command
	_, err = dg.ApplicationCommandCreate(
		dg.State.User.ID,
		dg.State.Application.GuildID,
		&discordgo.ApplicationCommand{
			Name:        "subscribe",
			Description: "Subscribe to an RSS feed",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "https://example.com/index.xml",
					Required:    true,
				},
			},
		},
	)
	if err != nil {
		slog.Error(fmt.Sprintf("error creating 'subscribe' command: %v", err))
		return
	}

	// add handler
	dg.AddHandler(dh.Create)

	// add event
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go dh.CheckNewEntries(ctx)

	// Set the playing status.
	_ = dg.UpdateGameStatus(0, "/subscribe <URL>")

	// Wait here until CTRL+C or other term signal is received
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
}
