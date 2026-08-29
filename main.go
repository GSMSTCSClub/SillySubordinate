package main

import (
	"log"

	"github.com/GSMSTCSClub/SillySubordinate/internal/bot"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading .env file")
	}

	bot.Start()
}
