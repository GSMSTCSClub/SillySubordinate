package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// The URL for the quote API to use. Currently ZenQuotes
const QuoteAPI string = "https://zenquotes.io/api/random"

// A struct to represent the ZenQuote JSON response.
type ZenQuote struct {
	Quote       string `json:"q"`
	Author      string `json:"a"`
	HTMLContent string `json:"h"`
}

var EightBallResponses []string = []string{
	"It is certain.",
	"It is decidedly so.",
	"Without a doubt.",
	"Yes, definitely.",
	"You may rely on it.",
	"As I see it, yes.",
	"Most likely.",
	"Outlook good.",
	"Yes.",
	"Signs point to yes.",
	"Reply hazy, try again.",
	"Ask again later.",
	"Better not tell you now.",
	"Cannot predict now.",
	"Concentrate and ask again.",
	"Don't count on it.",
	"My reply is no.",
	"My sources say no.",
	"Outlook not so good.",
	"Very doubtful.",
}

// A function to get a ZenQuote. Returns the ZenQuote struct.
func GetZenQuote() (*ZenQuote, error) {
	response, err := http.Get(QuoteAPI)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var quotes []ZenQuote
	err = json.NewDecoder(response.Body).Decode(&quotes)
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("No quotes fetched")
	}

	return &quotes[0], nil
}

// A function to turn float64 into pointers. For some aspects of discordgo.
func Float64Pointer(f float64) *float64 {
	return &f
}

// A function to turn int into pointers. For some aspects of discordgo.
func IntPointer(i int) *int {
	return &i
}

// Generate a random color integer. Can be used for embeds.
func RandomColor() int {
	// Generate random seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random number for color
	randomColor := rng.Intn(16777216)
	return randomColor
}
