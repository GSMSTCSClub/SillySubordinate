package util

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
)

// The URL for the quote API to use. Currently ZenQuotes
const QuoteAPI string = "https://zenquotes.io/api/random"

// Represents a quote from ZenQuotes. Contains the quote text, author, and HTML content from the JSON.
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

// A function that returns a [ZenQuote]
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

// Generate a random color [int]. Can be used for embeds.
func RandomColor() int {
	// Generate random number for color
	randomColor := rand.N(16777216)
	return randomColor
}

// Returns a random element in slice s. If slice s has a length of 0, the second return will return false
func RandomItem[E any](s []E) (E, bool) {
	if len(s) == 0 {
		return *new(E), false
	}

	randomIndex := rand.N(len(s))
	item := s[randomIndex]

	return item, true
}
