package router

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-shimada/discord-rss-bot/domain/model"
)

type discordHandler interface {
	Create(ds *discordgo.Session, dig *discordgo.InteractionCreate)
	List(ds *discordgo.Session, dig *discordgo.InteractionCreate)
	Delete(ds *discordgo.Session, dig *discordgo.InteractionCreate)
	FindAll() ([]model.Subscription, error)
	CheckNewEntries(ctx context.Context)
}

func NewRouter(token string) (*discordgo.Session, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return dg, nil
}

func Open(dg *discordgo.Session, dh discordHandler) {
	err := dg.Open()
	if err != nil {
		slog.Error(fmt.Sprintf("error opening connection: %v", err))
		return
	}
	defer dg.Close()

	// add subscribe command
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
	// add list command
	_, err = dg.ApplicationCommandCreate(
		dg.State.User.ID,
		dg.State.Application.GuildID,
		&discordgo.ApplicationCommand{
			Name:        "list",
			Description: "List all subscribed RSS feeds",
		},
	)
	if err != nil {
		slog.Error(fmt.Sprintf("error creating 'list' command: %v", err))
		return
	}
	// add unsubscribe command
	_, err = dg.ApplicationCommandCreate(
		dg.State.User.ID,
		dg.State.Application.GuildID,
		&discordgo.ApplicationCommand{
			Name:        "unsubscribe",
			Description: "Unsubscribe from an RSS feed",
			Options: []*discordgo.ApplicationCommandOption{
				{
					// Type:        discordgo.ApplicationCommandOptionString,
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "0",
					Required:    true,
				},
			},
		},
	)
	if err != nil {
		slog.Error(fmt.Sprintf("error creating 'list' command: %v", err))
		return
	}

	// add handler
	commandHandlers := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"subscribe":   dh.Create,
		"list":        dh.List,
		"unsubscribe": dh.Delete,
	}
	dg.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		},
	)

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
