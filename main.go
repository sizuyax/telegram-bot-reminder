package main

import (
	"fmt"
	"log"
	"psos/database"
	"psos/telegram"
)

func main() {
	_, err := database.InitDB()
	if err != nil {
		fmt.Println("failed to start database")
	}

	log.Println("bot started...")
	err2 := telegram.NewBotTelegram()
	if err2 != nil {
		fmt.Println("failed to start bot")
	}
}
