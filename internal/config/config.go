package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var (
	BotToken     string
	GuildID      string
	TestingMode  string
	CatchChannel string
	CatchItems   []string
	DataFile     string
	GameData     map[string]*User
)

type User struct {
	Languages []Language `json:"languages"`
}

type Language struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func ConfigVariables() {
	// Initialize bot constants
	BotToken = os.Getenv("BOT_TOKEN")
	GuildID = os.Getenv("GUILD_ID")
	TestingMode = os.Getenv("TESTING_MODE")
	CatchChannel = os.Getenv("CATCH_CHANNEL")
	CatchItems = []string{ // temporary items
		"Python",
		"JavaScript",
		"Go",
		"Gurt-Thon",
	}
	DataFile = os.Getenv("DATA_FILE")
}

func AccessUser(userID string) *User {
	user, userExists := GameData[userID]
	if !userExists {
		GameData[userID] = &User{
			Languages: []Language{},
		}
	}
	return user
}

func UpdateUser(userID string, userInfo *User) {
	GameData[userID] = userInfo
}

func LoadData() error {
	file, err := os.Open(DataFile)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		return err
	}

	// Make sure the file is closed later
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&GameData); err != nil {
		log.Printf("Failed to decode JSON file: %v", err)
		return err
	}

	return nil

}

func SaveData() error {
	jsonData, err := json.Marshal(GameData)
	if err != nil {
		log.Printf("Failed to save data: %v", err)
		return err
	}

	err = os.WriteFile(DataFile, jsonData, 0644)
	return nil
}
