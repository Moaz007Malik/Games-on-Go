package main

import (
	"fmt"
	"math/rand"
	"time"
)

var dictionary = []string{
	"Zombie",
	"Gopher",
	"Golang",
	"Games",
	"Hangman",
	"Pakistan",
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
}

func main() {
	rand.Seed(time.Now().UnixNano())

	targetWord := getRandomeWord()
	fmt.Println(targetWord)
}

func getRandomeWord() string {

	targetWord := dictionary[rand.Intn(len(dictionary))]
	return targetWord
}
