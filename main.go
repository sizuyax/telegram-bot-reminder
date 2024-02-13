package main

import (
	"fmt"
	"log"
	"psos/database"
	"psos/telegram"
)

func main() {
	log.Println("запуск базы данных")
	_, err := database.InitDB()
	if err != nil {
		fmt.Println("ошибка в main")
	}

	log.Println("запуск бота")
	_, err = telegram.NewBotTelegram()
	if err != nil {
		fmt.Println("ошибка в main")
	}
}
