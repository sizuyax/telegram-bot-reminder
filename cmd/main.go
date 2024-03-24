package main

import (
	"github.com/sirupsen/logrus"
	"remider/bot"
)

func main() {
	if err := bot.NewBotTelegram(); err != nil {
		logrus.Fatalf("failed to start bot, error: %v", err)
	}
}
