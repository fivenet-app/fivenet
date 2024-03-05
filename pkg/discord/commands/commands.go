package commands

import (
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/config"
	"go.uber.org/zap"
)

const GlobalCommandGuildID = "-1"

type CommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate)

type CommandFactory = func(cfg *config.BaseConfig) (*discordgo.ApplicationCommand, CommandHandler, error)

var (
	CommandsFactories = map[string]CommandFactory{}
	Commands          = map[string]*discordgo.ApplicationCommand{}
	CommandHandlers   = map[string]CommandHandler{}
)

type Cmds struct {
	logger *zap.Logger

	discord *discordgo.Session
}

func New(logger *zap.Logger, s *discordgo.Session, cfg *config.BaseConfig) (*Cmds, error) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	c := &Cmds{
		logger:  logger,
		discord: s,
	}

	for _, factory := range CommandsFactories {
		command, handler, err := factory(cfg)
		if err != nil {
			return nil, err
		}

		Commands[command.Name] = command
		CommandHandlers[command.Name] = handler
	}

	return c, nil
}

func (c *Cmds) RegisterGlobalCommands() error {
	cmds, err := c.discord.ApplicationCommands(c.discord.State.User.ID, "")
	if err != nil {
		return err
	}

	toRegister := []*discordgo.ApplicationCommand{}
	for _, command := range Commands {
		if command.GuildID == GlobalCommandGuildID {
			toRegister = append(toRegister, command)
		}
	}

	c.logger.Info("registering global commands", zap.Int("count", len(toRegister)))
	for _, command := range toRegister {
		if slices.ContainsFunc(cmds, func(cmd *discordgo.ApplicationCommand) bool {
			return cmd.Name == command.Name
		}) {
			c.logger.Debug(fmt.Sprintf("command '%v' already registered", command.Name))
			continue
		}

		if _, err := c.discord.ApplicationCommandCreate(c.discord.State.User.ID, "", command); err != nil {
			return fmt.Errorf("cannot create '%v' global command. %w", command.Name, err)
		}
	}

	return nil
}

func (c *Cmds) RemoveGuildCommands(guildID string) error {
	cmds, err := c.discord.ApplicationCommands(c.discord.State.User.ID, guildID)
	if err != nil {
		return err
	}

	// Remove guild registered commands
	c.logger.Debug("removing guild registered commands", zap.Any("commands", cmds))
	for _, command := range cmds {
		if err := c.discord.ApplicationCommandDelete(c.discord.State.User.ID, guildID, command.ID); err != nil {
			return fmt.Errorf("cannot delete '%v' guild registered command for guild '%s'. %w", command.Name, guildID, err)
		}
	}

	return nil
}

func GetDuplicateCommands(in []*discordgo.ApplicationCommand) []*discordgo.ApplicationCommand {
	allKeys := make(map[string]bool)
	list := []*discordgo.ApplicationCommand{}

	for _, item := range in {
		if _, value := allKeys[item.Name]; !value {
			allKeys[item.Name] = true
		} else {
			list = append(list, item)
		}
	}

	return list
}
