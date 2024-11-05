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
	"WeChat",
	"Pakistan Zindabad",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	targetWord := getRandomeWord()
	guessedLetters := initializeGuessedLetters(targetWord)
	hangmanState := 0

	for !isGameOver(targetWord, guessedLetters, hangmanState) {
		printGameState(targetWord, guessedLetters, hangmanState)
		input := readInput()
		if len(input) != 1 {
			fmt.Println("Invalid input. Please enter a single letter.")
			continue
		}

		letter := unicode.ToLower(rune(input[0]))

		if isCorrectGuess(targetWord, letter) {
			guessedLetters[letter] = true
		} else {
			hangmanState++
		}
	}

	printGameState(targetWord, guessedLetters, hangmanState)
	fmt.Println("Game Over... ")

	if isWordGuessed(targetWord, guessedLetters) {
		fmt.Println("You Win!")
	} else if isHangmanComplete(hangmanState) {
		fmt.Println("You Lose...")
	} else {
		panic("invalid state. Game is over and there is no winner!")
	}
}

func getRandomeWord() string {
	return dictionary[rand.Intn(len(dictionary))]
}

func initializeGuessedLetters(targetWord string) map[rune]bool {
	guessedLetters := map[rune]bool{}
	guessedLetters[unicode.ToLower(rune(targetWord[0]))] = true
	guessedLetters[unicode.ToLower(rune(targetWord[len(targetWord)-1]))] = true

	return guessedLetters
}

func isGameOver(targetWord string, guessedLetters map[rune]bool, hangmanState int) bool {
	return isWordGuessed(targetWord, guessedLetters) || isHangmanComplete(hangmanState)
}

func isWordGuessed(targetWord string, guessedLetters map[rune]bool) bool {
	for _, ch := range targetWord {
		if ch != ' ' && !guessedLetters[unicode.ToLower(ch)] {
			return false
		}
	}
	return true
}

func isHangmanComplete(hangmanState int) bool {
	return hangmanState >= 9
}

func printGameState(targetWord string, guessedLetters map[rune]bool, hangmanState int) {
	fmt.Println(getWordGuessingProgress(targetWord, guessedLetters))
	fmt.Println()
	fmt.Println(getHangmanDrawing(hangmanState))
}

func getWordGuessingProgress(targetWord string, guessedLetters map[rune]bool) string {
	var result strings.Builder
	for _, ch := range targetWord {
		if ch == ' ' {
			result.WriteString("  ")
		} else if guessedLetters[unicode.ToLower(ch)] {
			result.WriteRune(ch)
		} else {
			result.WriteString("_ ")
		}
	}
	return result.String()
}

func getHangmanDrawing(hangmanState int) string {
	data, err := ioutil.ReadFile(fmt.Sprintf("states/hangman%d", hangmanState))
	if err != nil {
		return fmt.Sprintf("Error loading drawing for state %d", hangmanState)
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

func isCorrectGuess(targetWord string, letter rune) bool {
	return strings.ContainsRune(strings.ToLower(targetWord), letter)
}
