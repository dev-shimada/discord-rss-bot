package discord

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/olekukonko/tablewriter"
)

type rssEntriesUsecase interface {
	Check(s model.Subscription) model.RssEntry
	CheckNewEntries(s []model.Subscription) []model.RssEntry
}

type subscriptionUsecase interface {
	FindAll() ([]model.Subscription, error)
	Create(sub model.Subscription) string
	Delete(sub model.Subscription) error
	List(sub model.Subscription) ([]model.Subscription, error)
}

type DiscordHandler struct {
	ds *discordgo.Session
	su subscriptionUsecase
	ru rssEntriesUsecase
}

func NewDiscordHandler(ds *discordgo.Session, su subscriptionUsecase, ru rssEntriesUsecase) DiscordHandler {
	return DiscordHandler{ds: ds, su: su, ru: ru}
}

func (d DiscordHandler) Create(ds *discordgo.Session, dic *discordgo.InteractionCreate) {
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

func (d DiscordHandler) List(ds *discordgo.Session, dic *discordgo.InteractionCreate) {
	// subscribe
	values, _ := d.su.List(model.Subscription{ChannelID: dic.ChannelID})

	title := "**Subscribed RSS feeds**"

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.Header([]string{"ID", "RSS URL"})
	for _, value := range values {
		if err := table.Append([]string{strconv.Itoa(int(value.ID)), value.RSSURL}); err != nil {
			slog.Error(fmt.Sprintf("Failed to append to table: %v", err))
		}
	}
	if err := table.Render(); err != nil {
		slog.Error(fmt.Sprintf("Failed to render table: %v", err))
	}

	txt := fmt.Sprintf("%s\n`%s`", title, tableString.String())

	_ = ds.InteractionRespond(dic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: txt,
		},
	})
}

func (d DiscordHandler) Delete(ds *discordgo.Session, dic *discordgo.InteractionCreate) {
	// get options
	options := dic.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}
	value := optionMap["id"].UintValue()

	// subscribe
	err := d.su.Delete(model.Subscription{ID: uint(value), ChannelID: dic.ChannelID})
	if err != nil {
		_ = ds.InteractionRespond(dic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to delete subscription.",
			},
		})
		return
	}

	_ = ds.InteractionRespond(dic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Successfully deleted subscription.",
		},
	})
}

func (d DiscordHandler) Check(ds *discordgo.Session, dic *discordgo.InteractionCreate) {
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
	rss := d.ru.Check(model.Subscription{RSSURL: rssUrl})
	if rss.EntryTitle == "" {
		_ = ds.InteractionRespond(dic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No new entries.",
			},
		})
		return
	}
	msg := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       rss.EntryTitle,
			URL:         rss.EntryLink,
			Description: rss.EntryTitle,
			Timestamp:   rss.PublishedAt.Format("2006-01-02 15:04:05"),
		},
	}
	if _, err := d.ds.ChannelMessageSendComplex(dic.ChannelID, msg); err != nil {
		slog.Error(fmt.Sprintf("Failed to send message: %v", err))
	}
	_ = ds.InteractionRespond(dic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "New entry found.",
		},
	})
}

func (d DiscordHandler) CheckNewEntries(ctx context.Context) {
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
