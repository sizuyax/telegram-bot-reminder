package main

import (
	"log"
	"psos/database"
	"psos/telegram"
)

func main() {
	_, err := database.InitDB()
	if err != nil {
		log.Fatalf("failed to connect to database, error: %v", err)
	}

	log.Println("bot started...")
	err = telegram.NewBotTelegram()
	if err != nil {
		log.Fatalf("failed to start bot, error: %v", err)
	}
}
