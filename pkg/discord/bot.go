package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/config"
	"go.uber.org/fx"
)

// TODO discord bot needs to join servers
// The bot will set user's roles when they have used the Discord third party login

type Bot struct {
	token string

	d *discordgo.Session
}

type BotParams struct {
	fx.In

	Config *config.Config
}

func NewBot(p BotParams) *Bot {
	return &Bot{
		token: p.Config.Discord.Bot.Token,
	}
}

func (b *Bot) Start() error {
	// Create a new Discord session using the provided login information.
	dg, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	b.d = dg

	return nil
}

func (b *Bot) SyncRoles() error {

	// TODO

	return nil
}
