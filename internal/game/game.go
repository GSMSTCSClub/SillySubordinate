package game

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/GSMSTCSClub/SillySubordinate/internal/config"
	"github.com/GSMSTCSClub/SillySubordinate/internal/util"
	"github.com/bwmarrin/discordgo"
)

// Holds data about the language catching game
type Game struct {
	CurrentSpawn   string
	SpawnActive    bool
	TargetChannel  string
	AvailableItems []string
}

func NewGame(targetChannel string, availableItems []string) *Game {
	game := &Game{
		CurrentSpawn:   "",
		SpawnActive:    false,
		TargetChannel:  targetChannel,
		AvailableItems: availableItems,
	}

	return game
}

func (g *Game) InitGame(s *discordgo.Session) {
	// Initialize ticker clock
	clock := time.NewTicker(30 * time.Second)

	for range clock.C {
		// Check if spawns are active
		if g.SpawnActive {
			continue
		}

		// Select random item to spawn
		var valid bool
		g.CurrentSpawn, valid = util.RandomItem(g.AvailableItems)
		if !valid {
			// Return if no items are in items list
			log.Fatal("Couldn't spawn items")
			return
		}
		g.SpawnActive = true

		s.ChannelMessageSend(g.TargetChannel, g.CurrentSpawn)
	}
}

func (g *Game) AttemptCatch(s *discordgo.Session, m *discordgo.MessageCreate, guess string) {
	if !g.SpawnActive {
		s.ChannelMessageSend(m.ChannelID, "No languages to catch!")
		return
	}

	if strings.EqualFold(guess, g.CurrentSpawn) {
		catcherID := m.Author.ID
		catchTime, err := discordgo.SnowflakeTimestamp(m.ID)
		if err != nil {
			log.Printf("Failed to get catch time: %v", err)
			return
		}

		language := &config.Language{
			Name:      g.CurrentSpawn,
			Timestamp: catchTime,
		}

		g.SpawnActive = false
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> caught **%s**!", catcherID, g.CurrentSpawn))
		// TODO: Database system
		catcher := config.AccessUser(catcherID)
		catcher.Languages = append(catcher.Languages, *language)
		config.UpdateUser(catcherID, catcher)
		err = config.SaveData()
		if err != nil {
			log.Printf("Failed to save data: %v", err)
			return
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Wrong name. Try again.")
	}
}
