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

	dg.AddHandler(dh.Create)

	// ticker := time.NewTicker(10 * time.Second)
	// defer ticker.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go dh.CheckNewEntries(ctx)

	// Wait here until CTRL+C or other term signal is received
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()

}
