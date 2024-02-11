package telegram

import (
	"fmt"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"psos/database"
	"time"
)

const (
	ModeHTML tb.ParseMode = "HTML"
)

var (
	menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	selector = &tb.ReplyMarkup{}

	btnState = menu.Text("Напомни мне статью!")
	btnEvent = menu.Text("Напомни мне событие!")
)

func NewBotTelegram() {
	fmt.Println("bot is running...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	tokenBot := os.Getenv("TG_API_TOKEN")

	b, err := tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
		Token:  tokenBot,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	handlerMessage(b)

	b.Start()
}

func handlerMessage(b *tb.Bot) {

	menu.Reply(
		menu.Row(btnState, btnEvent),
	)

	selector.Reply(
		selector.Row(btnState, btnEvent),
	)

	b.Handle("/start", func(m *tb.Message) {
		firstname := m.Sender.FirstName

		b.Send(m.Sender, "Привет, <b>"+firstname+"</b> этот бот поможет тебе напомнить самое важное! :)\n\n"+
			"Выбери что ты хочешь себе напомнить и когда!", menu, ModeHTML)
	})

	buttonHandler(b)
}

func buttonHandler(b *tb.Bot) {
	b.Handle(&btnState, func(m *tb.Message) {
		b.Send(m.Sender, "Отправь мне ссылку на ту статью которую хотел бы прочитать позже: ")
	})

	b.Handle(&btnEvent, func(m *tb.Message) {
		b.Send(m.Sender, "Напиши какое событие ты хотел бы напомнить себе позже: ")
		b.Handle(tb.OnText, func(m *tb.Message) {
			EventHandler(b, m)
		})
	})
}

func EventHandler(b *tb.Bot, m *tb.Message) {
	firstname := m.Sender.FirstName
	user_id := m.Sender.ID
	date := m.Time()
	text := m.Text

	event := database.Event{
		Firstname: firstname,
		Date:      date,
		User_id:   user_id,
		Text:      text,
	}

	err := database.CreateEvent(&event)
	if err != nil {
		log.Fatal(err)
	}
	b.Send(m.Sender, "Событие успешно добавлено в базу данных! :)")
}
