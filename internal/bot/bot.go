package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GSMSTCSClub/SillySubordinate/internal/commands"
	"github.com/bwmarrin/discordgo"
)

func Start() {
	// Initialize bot constants
	BotToken := os.Getenv("BOT_TOKEN")
	GuildID := os.Getenv("GUILD_ID")
	TestingMode := os.Getenv("TESTING_MODE")

	// Initialize bot dg
	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	/* session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore message if author is the bot
		if m.Author.ID == s.State.User.ID {
			return
		}

		args := strings.Split(m.Content, " ")
		if args[0] != Prefix {
			return
		}

		if args[1] == "hello" {
			s.ChannelMessageSendReply(m.ChannelID, "Hello, World!", m.Reference())
		}
	}) */

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	// Loop to register all definitions with Discord
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(commands.All))
	fmt.Println("Adding commands...")

	for _, cmd := range commands.All {
		// For testing, add guild ID
		m, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, cmd.Definition)
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

	// Wait until CTRL+C is pressed
	fmt.Println("Bot online! Logged in as " + dg.State.User.Username + "#" + dg.State.User.Discriminator)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Delete commands after shutting down
	if TestingMode == "YES" {
		log.Println("Removing commands...")
		for _, cmd := range registeredCommands {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, "", cmd.ID)
			if err != nil {
				log.Printf("Cannot delete '%v' command: %v", cmd.Name, err)
			}
		}
	}
}
