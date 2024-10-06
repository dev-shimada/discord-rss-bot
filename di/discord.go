package di

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/persistence"
	"github.com/dev-shimada/discord-rss-bot/interface/discord"
	"github.com/dev-shimada/discord-rss-bot/usecase"
	"gorm.io/gorm"
)

func DiscordHandler(db *gorm.DB, ds *discordgo.Session) discord.DiscordHandler {
	sr := persistence.NewSubscriptionPersistence(db)
	rr := persistence.NewRssEntryPersistence(db)
	su := usecase.NewSubscriptionUsecase(sr)
	ru := usecase.NewRssEntriesUsecase(rr)
	dh := discord.NewDiscordHandler(ds, su, ru)
	return dh
}
