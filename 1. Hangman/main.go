package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var inputReader = bufio.NewReader(os.Stdin)

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
	for {
		printGameState(targetWord, guessedLetters, hangmanState)
		input := readInput()
		if len(input) != 1 {
			fmt.Println("Invalid input. Please enter a single letter.")
			continue
		}
	}

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

func readInput() string {
	fmt.Print("> ")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(input)
}
