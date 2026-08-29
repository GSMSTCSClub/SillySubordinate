package commands

import (
	"github.com/bwmarrin/discordgo"
)

// TODO: Get user's languages and display them in an embed
func languages(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	msg := "Ill add this later"
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
}

func init() {
	register(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "languages",
			Description: "Lists a user's langauges!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Description: "The user who you would like to view languages",
					Required:    false, // If not given, get information on the interaction user
					Name:        "user",
				},
			},
		},
		Handler: languages,
	})
}
