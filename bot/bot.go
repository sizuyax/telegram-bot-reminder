package bot

import (
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"remider/config"
	"remider/database"
	"remider/handlers"
	"remider/repository"
	"remider/services"
	"time"
)

func NewBotTelegram() error {
	logrus.Println("bot is running...")

	config.LoadEnv()

	if err := database.InitDB(); err != nil {
		logrus.Fatalf("failed to connect to database, error: %v", err)
	}

	bot, err := tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
		Token:  config.Cfg.TgBotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		logrus.Fatal(err)
	}

	handlers.AllHandlers(bot)
	go SendRemindersToUsers(bot)

	bot.Start()
	return nil
}

func SendRemindersToUsers(b *tb.Bot) {

	log.Println("запуск отправки напоминаний")
	for {
		time.Sleep(1 * time.Second)
		repository.CheckEvents(b, services.Timezone.Timezone)
	}
}
