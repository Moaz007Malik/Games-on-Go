package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
	"unicode"
)

var dictionary = []string{
	"Zombie",
	"Gopher",
	"Golang",
	"Games",
	"Hangman",
	"Apple",
	"Microsoft",
	"Google",
	"Facebook",
	"Amazon",
	"Netflix",
	"YouTube",
	"Twitter",
	"Instagram",
	"WhatsApp",
	"Snapchat",
	"Tiktok",
	"LinkedIn",
	"Pinterest",
	"Reddit",
	"Tumblr",
	"Dropbox",
	"Github",
	"Gmail",
	"Skype",
	"Telegram",
	"Viber",
	"Whatsapp",
	"WeChat",
	"Pakistan Zindabad",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	targetWord := getRandomeWord()
	guessedLetters := initializeGuessedLetters(targetWord)
	hangmanState := 0
	printGameState(targetWord, guessedLetters, hangmanState)
}

func initializeGuessedLetters(targetWord string) map[rune]bool {
	guessedLetters := map[rune]bool{}
	guessedLetters[unicode.ToLower(rune(targetWord[0]))] = true
	guessedLetters[unicode.ToLower(rune(targetWord[len(targetWord)-1]))] = true

	return guessedLetters
}

func getRandomeWord() string {
	targetWord := dictionary[rand.Intn(len(dictionary))]
	return targetWord
}

func printGameState(
	targetWord string,
	guessedLetters map[rune]bool,
	hangmanState int,
) {
	fmt.Println(getWordGuessingProgess(targetWord, guessedLetters))
	fmt.Println()
	fmt.Println(getHangmanDrawing(hangmanState))
}

func getWordGuessingProgess(
	targetWord string,
	guessedLetters map[rune]bool,
) string {
	result := ""
	for _, ch := range targetWord {
		if ch == ' ' {
			result += " "
		} else if guessedLetters[unicode.ToLower(ch)] == true {
			result += fmt.Sprintf("%c", ch)
		} else {
			result += "_"
		}
		result += " "
	}
	return result
}

func getHangmanDrawing(hangmanState int) string {
	data, err := ioutil.ReadFile(fmt.Sprintf("states/hangman%d", hangmanState))
	if err != nil {
		panic(err)
	}

	return string(data)
}
