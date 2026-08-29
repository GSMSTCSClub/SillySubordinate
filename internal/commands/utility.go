package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/GSMSTCSClub/SillySubordinate/internal/util"
	"github.com/bwmarrin/discordgo"
)

func ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer ping response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Failed to defer interaction: %v", err)
		return
	}

	// Fetch response message
	responseMessage, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Printf("Failed to fetch message response: %v", err)
		return
	}

	// Get interaction creation timesamp
	interactionCreationTime, err := discordgo.SnowflakeTimestamp(i.ID)
	if err != nil {
		log.Printf("Failed to get interaction time: %v", err)
		return
	}

	// Get response creation timestamp
	responseCreationTime, err := discordgo.SnowflakeTimestamp(responseMessage.ID)
	if err != nil {
		log.Printf("Failed to get response time: %v", err)
		return
	}

	timeDifference := responseCreationTime.UnixMilli() - interactionCreationTime.UnixMilli()
	response := fmt.Sprintf("PONG!!!! :ping_pong: Took %d ms", timeDifference)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{

		Content: &response,
	})
}

func dumbPing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "PONG!!!! :ping_pong:",
		},
	})
	if err != nil {
		log.Printf("Failed to respond: %v", err)
	}
}

func purge(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Initialize response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	// Get interaction information
	channel := i.ChannelID
	messageCount := i.ApplicationCommandData().GetOption("count")
	interactionMember := i.Member

	// Check if user has adqueate permissions
	hasManageMessagePermission := (interactionMember.Permissions & discordgo.PermissionManageMessages) != 0

	if hasManageMessagePermission {
		messages, err := s.ChannelMessages(channel, int(messageCount.IntValue()), "", "", "")

		// Check for errors
		if err != nil {
			msg := "Couldn't fetch messages"
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
			return
		}

		// Fetch message IDs
		messageIDs := make([]string, len(messages))
		for idx, msg := range messages {
			messageIDs[idx] = msg.ID
		}

		// Bulk delete
		err = s.ChannelMessagesBulkDelete(channel, messageIDs)
		if err != nil {
			msg := "Failed to delete messages"
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
			return
		}

		msg := fmt.Sprintf("Deleted %d messages", messageCount.IntValue())
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
	} else {
		msg := "Sorry, you don't have the proper permissions"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
	}
}

func user(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Initialize response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	user := i.Member.User
	userOption := i.ApplicationCommandData().GetOption("user")
	if userOption != nil {
		user = userOption.UserValue(s)
	}

	userAvatar := user.AvatarURL("512")

	// Get user creation date
	creationDate, err := discordgo.SnowflakeTimestamp(user.ID)
	if err != nil {
		log.Printf("Couldn't get user creation timestamp: %v", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:     "Stats of " + user.Username,
		Color:     util.RandomColor(),
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: userAvatar},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Username",
				Value:  user.Username,
				Inline: true,
			},
			{
				Name:   "Display Name",
				Value:  user.DisplayName(),
				Inline: true,
			},
			{
				Name:   "ID",
				Value:  user.ID,
				Inline: true,
			},
			{
				Name:   "Discord Join Date",
				Value:  creationDate.UTC().Format(time.UnixDate),
				Inline: true,
			},
		},
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			embed,
		},
	})
}

func init() {
	register(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Test the bot's ping.",
		},
		Handler: ping,
	})

	register(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "dumb-ping",
			Description: "A dumb version of /ping",
		},
		Handler: dumbPing,
	})

	register(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "purge",
			Description: "Purge messages within a channel!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "count",
					Required:    true,
					Description: "Number of messages you would like to purge",
					MinValue:    new(1.0),
					MaxValue:    *new(100.0),
				},
			},
		},
		Handler: purge,
	})

	register(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "user",
			Description: "Gets information about a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Description: "The user who you would like to get information on",
					Required:    false, // If not given, get information on the interaction user
					Name:        "user",
				},
			},
		},
		Handler: user,
	})
}
