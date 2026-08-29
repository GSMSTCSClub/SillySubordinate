package commands

import (
	"math/rand/v2"

	"github.com/GSMSTCSClub/SillySubordinate/internal/util"
	"github.com/bwmarrin/discordgo"
)

func quote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	member := i.Interaction.Member

	// Initialize response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Fetch from ZenQuotes
	quote, err := util.GetZenQuote()
	if err != nil {
		msg := "Error getting quote"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		return
	}

	// Choose appropriate name for embed header
	var nameTitle string
	if member.Nick != "" {
		// If the user has a nickname, use that
		nameTitle = member.Nick
	} else {
		// Otherwise, use the user's display name
		nameTitle = member.DisplayName()
	}

	// Concatenate text for embed
	quoteText := "“" + quote.Quote + "”"
	quoteAuthor := "\\- " + quote.Author

	// Create embed to send to user
	embed := &discordgo.MessageEmbed{
		Title: nameTitle + "'s Quote",
		Color: util.RandomColor(),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  quoteText,
				Value: quoteAuthor,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quotes from https://zenquotes.io",
		},
	}

	// Response to interaction with embed
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			embed,
		},
	})
}

func eightBall(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get interaction options
	question := i.ApplicationCommandData().GetOption("question").StringValue()

	hiddenOption := i.ApplicationCommandData().GetOption("hidden")
	var hidden bool
	if hiddenOption != nil {
		hidden = hiddenOption.BoolValue()
	}

	// Initialize response
	if hidden {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	}

	// Generate random 8 ball response
	randomIndex := rand.N(len(util.EightBallResponses))
	randomResponse := util.EightBallResponses[randomIndex]

	// Create embed to send to user
	embed := &discordgo.MessageEmbed{
		Title: "Answer to " + question,
		Color: util.RandomColor(),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Your answer is.... :8ball:",
				Value:  randomResponse,
				Inline: false,
			},
		},
	}

	// Send the final response
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			embed,
		},
	})
}

func init() {
	register(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "quote",
			Description: "Get an inspirational quote",
		},
		Handler: quote,
	})

	register(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "8ball",
			Description: "Unsure on a decision? Ask the magic 8 ball your question...",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Required:    true,
					Description: "The question you want to ask the 8 ball",
					MinLength:   util.IntPointer(5),
					MaxLength:   *util.IntPointer(100),
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "hidden",
					Required:    false,
					Description: "Whether you would like others to see your question and response (false by default).",
				},
			},
		},
		Handler: eightBall,
	})
}
