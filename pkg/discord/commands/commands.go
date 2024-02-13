package commands

import (
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands        = []*discordgo.ApplicationCommand{}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

type Cmds struct {
	discord *discordgo.Session
}

func New(s *discordgo.Session) *Cmds {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	return &Cmds{
		discord: s,
	}
}

func (c *Cmds) Register(guild *discordgo.Guild) error {
	cmds, err := c.discord.ApplicationCommands(c.discord.State.User.ID, guild.ID)
	if err != nil {
		return err
	}

	for _, v := range Commands {
		if slices.ContainsFunc(cmds, func(cmd *discordgo.ApplicationCommand) bool {
			return cmd.Name == v.Name
		}) {
			continue
		}

		if _, err := c.discord.ApplicationCommandCreate(c.discord.State.User.ID, guild.ID, v); err != nil {
			return fmt.Errorf("cannot create '%v' command for guild '%s': %v", v.Name, guild.ID, err)
		}
	}

	return nil
}
