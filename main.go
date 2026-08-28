package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Hello, World!")

	// unnecessary stuff
	discord, err := discordgo.New("Bot" + "imagine a token")
	if err != nil {
		fmt.Println("couldnt make yo bot")
		return
	}

	discord.Close()
}
