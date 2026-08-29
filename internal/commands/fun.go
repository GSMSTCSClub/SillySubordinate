package commands

import (
	"github.com/bwmarrin/discordgo"
)

func quote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// stuff
}

func init() {
	register(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "quote",
			Description: "Get an inspirational quote",
		},
		Handler: quote,
	})
}
