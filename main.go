package main

import (
	"fmt"
	"psos/database"
	"psos/telegram"
)

func main() {
	err := database.InitDB()
	if err != nil {
		fmt.Println("ошибка в main")
	}
	telegram.NewBotTelegram()

}
