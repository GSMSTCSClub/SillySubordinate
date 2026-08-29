package bot

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/GSMSTCSClub/SillySubordinate/internal/commands"
	"github.com/GSMSTCSClub/SillySubordinate/internal/config"
	"github.com/GSMSTCSClub/SillySubordinate/internal/game"
	"github.com/bwmarrin/discordgo"
)

func Start() {
	config.ConfigVariables()

	// Initialize catching game
	catchGame := game.NewGame(config.CatchChannel, config.CatchItems)
	err := config.LoadData()
	if err != nil {
		log.Printf("Failed to load config data: %v", err)
		return
	}

	// Initialize bot dg
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Create goroutine for the catching game
	go catchGame.InitGame(dg)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	// Loop to register all definitions with Discord
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(commands.All))
	log.Println("Adding commands...")

	for _, cmd := range commands.All {
		// For testing, add guild ID
		m, err := dg.ApplicationCommandCreate(dg.State.User.ID, config.GuildID, cmd.Definition)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", cmd.Definition.Name, err)
		}
		registeredCommands = append(registeredCommands, m)
	}

	// Add route handler
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		// Look up command using name
		if cmd, exists := commands.All[i.ApplicationCommandData().Name]; exists {
			cmd.Handler(s, i)
		}
	})

	// Catch game handler
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, "!catch ") {
			guess := strings.TrimSpace(strings.TrimPrefix(m.Content, "!catch "))
			catchGame.AttemptCatch(s, m, guess)
		}
	})

	// Wait until CTRL+C is pressed
	log.Println("Bot online! Logged in as " + dg.State.User.Username + "#" + dg.State.User.Discriminator)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Delete commands after shutting down
	if config.TestingMode == "YES" {
		log.Println("Removing commands...")
		for _, cmd := range registeredCommands {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, config.GuildID, cmd.ID)
			if err != nil {
				log.Printf("Cannot delete '%v' command: %v", cmd.Name, err)
			}
		}
	}
}
