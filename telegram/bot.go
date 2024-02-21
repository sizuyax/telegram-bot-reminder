package telegram

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"os"
	"psos/database"
	"psos/utils"
	"strconv"
	"time"
)

const (
	ModeHTML tb.ParseMode = "HTML"
)

var (
	menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	selector = &tb.ReplyMarkup{}

	btnEvent     = menu.Text("Нагадай мені подію!")
	btnEditName  = menu.Text("Змінити назву нагадування!")
	btnCheck     = menu.Text("Перевірити нагадування.")
	btnEditDate  = menu.Text("Змінити дату нагадування.")
	btnDelete    = menu.Text("Видалити нагадування.")
	btnDeleteAll = menu.Text("Видалити всі нагадування.")
)

var (
	text      string
	user_id   int64
	firstname string
)

// start bot
func NewBotTelegram() error {
	fmt.Println("bot is running...")

	utils.LoadEnv()
	tokenBot := os.Getenv("TG_API_TOKEN")

	b, err := tb.NewBot(tb.Settings{
		URL:    "https://your-heroku-app.herokuapp.com",
		Token:  tokenBot,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	handlerMessage(b)

	go SendReminder(b)

	b.Handle("/"+b.Token, func(m *tb.Message) {
		b.ProcessUpdate(tb.Update{Message: m})
	})

	// Запуск веб-сервера на указанном порту
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // По умолчанию используем порт 8080, если переменная окружения не задана
	}

	go func() {
		log.Printf("Starting server on port %s...", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()

	b.Start()
	return nil
}

func handlerMessage(b *tb.Bot) {

	menu.Reply(
		menu.Row(btnEvent, btnCheck),
		menu.Row(btnEditName, btnEditDate),
		menu.Row(btnDelete, btnDeleteAll),
	)

	selector.Reply(
		selector.Row(btnEvent, btnEditName),
		selector.Row(btnEditName, btnEditDate),
		selector.Row(btnDelete, btnDeleteAll),
	)

	b.Handle("/start", func(m *tb.Message) {
		firstname := m.Sender.FirstName

		b.Send(m.Sender, "Привiт, <b>"+firstname+"</b> цей бот допоможе тобі нагадати найважливіше! :)\n\n"+
			"Вибери, що ти хочеш собі нагадати і коли!", menu, ModeHTML)
	})

	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Sender, "Ось список доступних команд:\n\n"+
			"/start - почати роботу з ботом\n\n"+
			"/help - список доступних команд\n\n"+
			"/check - перевірити всі нагадування\n\n"+
			"/event - додати нове нагадування\n\n"+
			"/editname - змінити текст нагадування\n\n"+
			"/editdate - змінити дату нагадування\n\n"+
			"/delete - видалити нагадування\n\n"+
			"/deleteall - видалити всі нагадування\n\n"+
			"/stop - завершити роботу з ботом", ModeHTML)
	})

	b.Handle("/check", func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)

		text := database.GetText()

		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування ")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
		}
	})

	b.Handle("/editname", func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)

		text := database.GetText()
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, яке потрібно змінити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				EdRem(b, m)
			})
		}
	})

	b.Handle("/event", func(m *tb.Message) {
		b.Send(m.Sender, "Напиши яку подію ти хотів би нагадати собі пізніше: ")
		b.Handle(tb.OnText, func(m *tb.Message) {

			EventText(m)

			b.Send(m.Sender, "Тепер напиши коли тобі нагадати твою подію у форматі 2021-12-31 23:59 часовий пояс UTC: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				EventHandler(b, m)
			})
		})
	})

	b.Handle("/editdate", func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)

		text := database.GetText()
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, а потiм дату якого потрібно змінити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				EditTimeForReminder(b, m)
			})
		}
	})

	b.Handle("/delete", func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)
		text := database.GetText()
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, яке потрібно видалити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				database.DeleteReminder(m.Text, b, m)
				b.Send(m.Sender, "Нагадування успішно видалено!")
			})
		}
	})

	b.Handle("/deleteall", func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)
		text := database.GetText()
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Ви впевнені, що хочете видалити всі нагадування? Введіть 'Так' або 'Ні': ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				if m.Text == "Так" {
					database.DeleteAllReminders()
					b.Send(m.Sender, "Всі нагадування успішно видалено!")
				} else if m.Text == "Ні" {
					b.Send(m.Sender, "Ви відмінили видалення нагадувань!")
				} else {
					b.Send(m.Sender, "Введіть 'Так' або 'Ні': ")
				}
			})
		}

	})

	b.Handle("/stop", func(m *tb.Message) {
		b.Send(m.Sender, "Бувай! Якщо тобі щось важливе раптом згадається, звертайся! :)")
	})

	buttonHandler(b)
}

func buttonHandler(b *tb.Bot) {

	allNoNeedsHandlers(b)

	b.Handle(&btnCheck, func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)
		//var reminderToUpdate *database.Event

		text := database.GetText()

		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування ")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
		}
	})

	b.Handle(&btnEditName, func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)
		texts := database.GetText()

		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування ")
			for _, text := range texts {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, який потрібно змінити:")
			b.Handle(tb.OnText, func(m *tb.Message) {
				EdRem(b, m)
			})
		}
	})

	b.Handle(&btnEvent, func(m *tb.Message) {
		b.Send(m.Sender, "Напиши яку подію ти хотів би нагадати собі пізніше: ")
		b.Handle(tb.OnText, func(m *tb.Message) {

			EventText(m)

			b.Send(m.Sender, "Тепер напиши коли тобі нагадати твою подію у форматі 2023-12-31 23:59 часовий пояс UTC:")
			b.Handle(tb.OnText, func(m *tb.Message) {
				EventHandler(b, m)
			})
		})
	})

	b.Handle(&btnEditDate, func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)
		text := database.GetText()

		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування ")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, а потiм дату якого потрібно змінити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				EditTimeForReminder(b, m)
			})
		}
	})

	b.Handle(&btnDelete, func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)
		text := database.GetText()
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, яке потрібно видалити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				database.DeleteReminder(m.Text, b, m)
				b.Send(m.Sender, "Нагадування успішно видалено!")
			})
		}
	})

	b.Handle(&btnDeleteAll, func(m *tb.Message) {
		reminders := database.GetReminders(m.Sender.ID)
		text := database.GetText()
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Ви впевнені, що хочете видалити всі нагадування? Введіть 'Так' або 'Ні': ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				if m.Text == "Так" {
					database.DeleteAllReminders()
					b.Send(m.Sender, "Всі нагадування успішно видалено!")
				} else if m.Text == "Ні" {
					b.Send(m.Sender, "Ви відмінили видалення нагадувань!")
				} else {
					b.Send(m.Sender, "Введіть 'Так' або 'Ні': ")
				}
			})
		}
	})
}

// block for events

func EventHandler(b *tb.Bot, m *tb.Message) {

	date, err := time.Parse("2006-01-02 15:04", m.Text)
	if err != nil {
		b.Send(m.Sender, "Неправильний формат дати, спробуй ще раз!")
		return
	}
	if date.Before(time.Now()) {
		b.Send(m.Sender, "Дата не може бути раніше за поточну, спробуй ще раз!")
		return
	}

	event := database.Event{
		Firstname: firstname,
		Date:      date,
		User_id:   user_id,
		Text:      text,
	}

	err = database.CreateEvent(&event, m)
	if err != nil {
		log.Fatal(err)
	}
	b.Send(m.Sender, "Подія успішно додана до бази даних! :)")
}

func EventText(m *tb.Message) {
	text = m.Text
	firstname = m.Sender.FirstName
	user_id = m.Sender.ID
}

func SendReminder(b *tb.Bot) {

	log.Println("запуск отправки напоминаний")
	for {
		time.Sleep(1 * time.Second)
		database.CheckReminders(b)
	}
}

func EditTimeForReminder(b *tb.Bot, m *tb.Message) {
	reminders := database.GetReminders(m.Sender.ID)

	reminderID, err := strconv.Atoi(m.Text)
	if err != nil {
		b.Send(m.Sender, "Неправильний формат нагадування!")
		return
	}

	foundReminder := findReminderByID(reminders, reminderID)

	// Если найдено напоминание с указанным ID
	if foundReminder != nil {
		// Запросить новую дату для напоминания
		b.Send(m.Sender, "Введіть нову дату нагадування 2023-12-31 23:59 часовий пояс UTC: ")

		b.Handle(tb.OnText, func(m *tb.Message) {
			updateReminderDate(b, m, foundReminder)
		})
	} else {
		// Не найдено напоминание с указанным ID
		b.Send(m.Sender, "У вас немає нагадування із таким номером!")
	}
}

func EdRem(b *tb.Bot, m *tb.Message) {
	reminders := database.GetReminders(m.Sender.ID)

	reminderID, err := strconv.Atoi(m.Text)
	if err != nil {
		b.Send(m.Sender, "Неправильний формат нагадування!")
		return
	}

	foundReminder := findReminderByID(reminders, reminderID)

	// Если найдено напоминание с указанным ID
	if foundReminder != nil {
		// Запросить новый текст напоминания
		b.Send(m.Sender, "Введіть новий текст для нагадування:")

		b.Handle(tb.OnText, func(m *tb.Message) {
			newText := m.Text

			// Обновить текст напоминания
			foundReminder.Text = newText

			// Обновить напоминание в базе данных
			err := database.UpdateReminder(foundReminder)
			if err != nil {
				log.Fatal(err)
			}

			b.Send(m.Sender, "Нагадування успішно змінено!")
		})
	} else {
		// Не найдено напоминание с указанным ID
		b.Send(m.Sender, "У вас немає нагадування із таким номером!")
	}

}

func findReminderByID(reminders []database.Event, reminderID int) *database.Event {
	for i := range reminders {
		if int(reminders[i].ID) == reminderID {
			return &reminders[i]
		}
	}
	return nil
}

func updateReminderDate(b *tb.Bot, m *tb.Message, foundReminder *database.Event) {
	date, err := time.Parse("2006-01-02 15:04", m.Text)
	if err != nil {
		b.Send(m.Sender, "Неправильний формат дати, спробуй ще раз!")
		return
	}
	if date.Before(time.Now()) {
		b.Send(m.Sender, "Дата не може бути раніше за поточну, спробуй ще раз!")
		return
	}

	// Обновить дату напоминания
	foundReminder.Date = date

	// Обновить напоминание в базе данных
	err = database.UpdateReminder(foundReminder)
	if err != nil {
		log.Fatal(err)
	}

	b.Send(m.Sender, "Дата нагадування успішно змінено!")
}

// end block for events

func allNoNeedsHandlers(b *tb.Bot) {
	b.Handle(tb.OnText, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnAudio, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnDocument, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnSticker, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnVideo, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnVoice, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnVideoNote, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnContact, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")

	})
	b.Handle(tb.OnLocation, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnVenue, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnPoll, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})
	b.Handle(tb.OnDice, func(m *tb.Message) {
		b.Send(m.Sender, "Не розумію тебе, вибери щось із меню або напиши /start")
	})

}
