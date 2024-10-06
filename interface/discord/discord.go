package discord

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/usecase"
)

type DiscordHandler interface {
	Create(ds *discordgo.Session, dm *discordgo.MessageCreate)
	FindAll() ([]model.Subscription, error)
	CheckNewEntries(ctx context.Context)
}

type discordHandler struct {
	ds *discordgo.Session
	su usecase.SubscriptionUsecase
	ru usecase.RssEntriesUsecase
}

func NewDiscordHandler(ds *discordgo.Session, su usecase.SubscriptionUsecase, ru usecase.RssEntriesUsecase) DiscordHandler {
	return &discordHandler{ds: ds, su: su, ru: ru}
}

func (d discordHandler) Create(ds *discordgo.Session, dm *discordgo.MessageCreate) {
	if dm.Author.ID == ds.State.User.ID {
		return
	}
	if !strings.HasPrefix(dm.Content, "!subscribe ") {
		return
	}

	// validate URL
	arr := strings.Split(dm.Content, " ")
	value := arr[1]
	rssURL, err := url.Parse(value)
	if err != nil {
		if _, err := ds.ChannelMessageSend(dm.ChannelID, "Invalid URL."); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
		return
	}

	// subscribe
	msg := d.su.Create(model.Subscription{ChannelID: dm.ChannelID, RSSURL: rssURL.String()})
	if _, err := ds.ChannelMessageSend(dm.ChannelID, msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (d discordHandler) FindAll() ([]model.Subscription, error) {
	return d.su.FindAll()
}

func (d discordHandler) CheckNewEntries(ctx context.Context) {
	t := time.NewTicker(10 * time.Minute)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			subs, err := d.su.FindAll()
			if err != nil {
				log.Fatalf("error fetching subscriptions: %v", err)
				return
			}
			newEntries := d.ru.CheckNewEntries(subs)
			for _, entry := range subs {
				for _, newEntry := range newEntries {
					if entry.RSSURL == newEntry.RSSURL {
						msg := &discordgo.MessageSend{
							Embed: &discordgo.MessageEmbed{
								Title:       newEntry.EntryTitle,
								URL:         newEntry.EntryLink,
								Description: newEntry.EntryTitle,
								Timestamp:   newEntry.PublishedAt.Format("2006-01-02 15:04:05"),
							},
						}
						if _, err := d.ds.ChannelMessageSendComplex(entry.ChannelID, msg); err != nil {
							log.Printf("Failed to send message: %v", err)
						}
					}
				}
			}
		}
	}
}
