package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/discord/modules"
	"go.uber.org/zap"
)

type Guild struct {
	Job string `alias:"job"`
	ID  string `alias:"id"`

	bot     *Bot
	guild   *discordgo.Guild
	modules map[string]modules.Module
}

func NewGuild(b *Bot, guild *discordgo.Guild, job string) (*Guild, error) {
	ms := map[string]modules.Module{}

	base := modules.NewBaseModule(b.ctx, b.logger, b.db, b.discord, guild, job, b.cfg, b.enricher)
	var err error
	for name := range modules.Modules {
		ms[name], err = modules.GetModule(name, base)
		if err != nil {
			return nil, err
		}
	}

	return &Guild{
		bot:     b,
		guild:   guild,
		ID:      guild.ID,
		Job:     job,
		modules: ms,
	}, nil
}

func (g *Guild) Run() error {
	for key, module := range g.modules {
		g.bot.logger.Debug("running discord guild module", zap.String("dc_module", key))
		if err := module.Run(); err != nil {
			return err
		}
	}

	return nil
}
