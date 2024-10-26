package discord

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/usecase"
)

type DiscordHandler interface {
	Create(ds *discordgo.Session, dig *discordgo.InteractionCreate)
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

func (d discordHandler) Create(ds *discordgo.Session, dic *discordgo.InteractionCreate) {
	// get options
	options := dic.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}
	value := optionMap["url"].StringValue()

	// validate URL
	validUrl, err := url.ParseRequestURI(value)
	if err != nil {
		_ = ds.InteractionRespond(dic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid URL.",
			},
		})
		return
	}
	rssUrl := validUrl.String()
	// subscribe
	d.su.Create(model.Subscription{ChannelID: dic.ChannelID, RSSURL: rssUrl})
	_ = ds.InteractionRespond(dic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Successfully subscribed to RSS feed: %s", rssUrl),
		},
	})
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
				slog.Warn(fmt.Sprintf("error fetching subscriptions: %v", err))
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
							slog.Error(fmt.Sprintf("Failed to send message: %v", err))
						}
					}
				}
			}
		}
	}
}
