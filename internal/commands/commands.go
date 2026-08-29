package commands

import "github.com/bwmarrin/discordgo"

// Define a Command struct
type Command struct {
	Definition *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Create variable to store commands
var All = make(map[string]Command)

func register(cmd Command) {
	All[cmd.Definition.Name] = cmd
}
