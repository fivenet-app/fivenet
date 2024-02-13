package commands

import (
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/utils"
	"go.uber.org/zap"
)

var (
	Commands        = []*discordgo.ApplicationCommand{}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

type Cmds struct {
	logger *zap.Logger

	discord *discordgo.Session
}

func New(logger *zap.Logger, s *discordgo.Session) *Cmds {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	return &Cmds{
		logger:  logger,
		discord: s,
	}
}

func (c *Cmds) Register(guild *discordgo.Guild) error {
	cmds, err := c.discord.ApplicationCommands(c.discord.State.User.ID, guild.ID)
	if err != nil {
		return err
	}

	c.logger.Info("registering commands", zap.Int("count", len(Commands)))
	for _, command := range Commands {
		if slices.ContainsFunc(cmds, func(cmd *discordgo.ApplicationCommand) bool {
			return cmd.Name == command.Name
		}) {
			continue
		}

		if _, err := c.discord.ApplicationCommandCreate(c.discord.State.User.ID, guild.ID, command); err != nil {
			return fmt.Errorf("cannot create '%v' command for guild '%s': %v", command.Name, guild.ID, err)
		}
	}

	return nil
}

func (c *Cmds) Unregister(guild *discordgo.Guild) error {
	cmds, err := c.discord.ApplicationCommands(c.discord.State.User.ID, guild.ID)
	if err != nil {
		return err
	}

	// Remove duplicate commands
	duplicates := utils.GetDuplicates(cmds)
	c.logger.Info("removing duplicate commands", zap.Int("duplicates", len(duplicates)))
	for _, command := range duplicates {
		if err := c.discord.ApplicationCommandDelete(c.discord.State.User.ID, guild.ID, command.ID); err != nil {
			return fmt.Errorf("cannot delete '%v' command for guild '%s': %v", command.Name, guild.ID, err)
		}
	}

	return nil
}
