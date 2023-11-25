package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands        = []*discordgo.ApplicationCommand{}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

type Cmds struct {
	guildId            string
	RegisteredCommands []*discordgo.ApplicationCommand
}

func New(guildId string) *Cmds {
	return &Cmds{
		guildId:            guildId,
		RegisteredCommands: []*discordgo.ApplicationCommand{},
	}
}

func (c *Cmds) Register(s *discordgo.Session) error {
	for _, v := range Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, c.guildId, v)
		if err != nil {
			return fmt.Errorf("cannot create '%v' command for guild '%s': %v", v.Name, c.guildId, err)
		}
		c.RegisteredCommands = append(c.RegisteredCommands, cmd)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	return nil
}
