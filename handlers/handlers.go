package handlers

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"remider/repository"
	"remider/services"
)

const (
	ModeHTML tb.ParseMode = "HTML"
)

var (
	menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	selector = &tb.ReplyMarkup{}

	btnEvent     = menu.Text("Нагадай мені подію!")
	btnEditName  = menu.Text("Змінити назву нагадування.")
	btnCheck     = menu.Text("Перевірити нагадування.")
	btnEditDate  = menu.Text("Змінити дату нагадування.")
	btnDelete    = menu.Text("Видалити нагадування.")
	btnDeleteAll = menu.Text("Видалити всі нагадування.")
)

func AllHandlers(b *tb.Bot) {

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

	buttonHandlers(b)
	commandHandlers(b)
	allNoNeedsHandlers(b)
}

func commandHandlers(b *tb.Bot) {
	b.Handle("/start", func(m *tb.Message) {
		firstname := m.Sender.FirstName

		b.Send(m.Sender, "Привiт, <b>"+firstname+"</b> цей бот допоможе тобі нагадати найважливіше! :)\n\n"+
			"Напиши своє мiсто щоб бот мiг виставити часовий пояс: ", menu, ModeHTML)
		b.Handle(tb.OnText, func(m *tb.Message) {
			services.GetCity(m.Text)
			b.Send(m.Sender, "Успiшно, тепер можете працювати з нагадуваннями!")
		})
	})

	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Sender, "Ось список доступних команд:\n\n"+
			"/start - почати роботу з ботом\n\n"+
			"/help - список доступних команд\n\n"+
			"/setcity - виставити свое мicто\n\n"+
			"/check - перевірити всі нагадування\n\n"+
			"/event - додати нове нагадування\n\n"+
			"/editname - змінити текст нагадування\n\n"+
			"/editdate - змінити дату нагадування\n\n"+
			"/delete - видалити нагадування\n\n"+
			"/deleteall - видалити всі нагадування\n\n"+
			"/stop - завершити роботу з ботом", ModeHTML)
	})

	b.Handle("/setcity", func(m *tb.Message) {
		b.Send(m.Sender, "Напишiть своє мicто: ")

		b.Handle(tb.OnText, func(m *tb.Message) {
			services.GetCity(m.Text)
			b.Send(m.Sender, "Мiсто успiшно встановлено!")
		})
	})

	b.Handle("/check", func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		texts := repository.GetText(reminders)
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування: ")
			for _, text := range texts {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
		}
	})

	b.Handle("/editname", func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		if reminders == nil {
			b.Send(m.Sender, "Ви не вписали своє мiсто для часового поясу!\n\n"+
				"Щоб вписати свое мiсто нажмiть на команду /setcity")
			return
		}

		texts := repository.GetText(reminders)
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування: ")
			for _, text := range texts {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, яке потрібно змінити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				services.EditTextForEvent(b, m)
			})
		}
	})

	b.Handle("/event", func(m *tb.Message) {
		b.Send(m.Sender, "Напиши яку подію ти хотів би нагадати собі пізніше: ")
		b.Handle(tb.OnText, func(m *tb.Message) {

			services.SaveEventText(m)

			b.Send(m.Sender, "Тепер напиши коли тобі нагадати твою подію у форматі <b>2021-12-31 23:59</b>: ", ModeHTML)
			b.Handle(tb.OnText, func(m *tb.Message) {
				services.SaveEvent(b, m)
			})
		})
	})

	b.Handle("/editdate", func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		if reminders == nil {
			b.Send(m.Sender, "Ви не вписали cвоє мiсто для часового поясу!\n\n"+
				"Щоб вписати своє мiсто нажмiть на команду /setcity")
			return
		}

		texts := repository.GetText(reminders)
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування: ")
			for _, text := range texts {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}

			b.Send(m.Sender, "Введіть номер нагадування, а потiм дату якого потрібно змінити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				services.EditTimeForEvent(b, m)
			})
		}
	})

	b.Handle("/delete", func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		text := repository.GetText(reminders)
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування: ")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, яке потрібно видалити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				repository.DeleteEvent(m.Text, b, m)
				b.Send(m.Sender, "Нагадування успішно видалено!")
			})
		}
	})

	b.Handle("/deleteall", func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		texts := repository.GetText(reminders)
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування: ")
			for _, text := range texts {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Ви впевнені, що хочете видалити всі нагадування? Введіть 'Так' або 'Ні': ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				if m.Text == "Так" {
					repository.DeleteAllEvents()
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
}

func buttonHandlers(b *tb.Bot) {
	b.Handle(&btnCheck, func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		texts := repository.GetText(reminders)

		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування ")
			for _, text := range texts {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
		}
	})

	b.Handle(&btnEditName, func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		if reminders == nil {
			b.Send(m.Sender, "Ви не вписали свое мiсто для вибору часового поясу!\n"+
				"Щоб вписати свое мiсто нажмiть на команду /setcity")
			return
		}

		texts := repository.GetText(reminders)
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range texts {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, яке потрібно змінити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				services.EditTextForEvent(b, m)
			})
		}
	})

	b.Handle(&btnEvent, func(m *tb.Message) {
		b.Send(m.Sender, "Напиши яку подію ти хотів би нагадати собі пізніше: ")

		b.Handle(tb.OnText, func(m *tb.Message) {

			services.SaveEventText(m)

			b.Send(m.Sender, "Тепер напиши коли тобі нагадати твою подію у форматі 2021-12-31 23:59: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				services.SaveEvent(b, m)
			})
		})
	})

	b.Handle(&btnEditDate, func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		if reminders == nil {
			b.Send(m.Sender, "Ви не вписали свое мiсто для вибору часового поясу!\n"+
				"Щоб вписати свое мiсто нажмiть на команду /setcity")
			return
		}

		texts := repository.GetText(reminders)
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range texts {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}

			b.Send(m.Sender, "Введіть номер нагадування, а потiм дату якого потрібно змінити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				services.EditTimeForEvent(b, m)
			})
		}
	})

	b.Handle(&btnDelete, func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		text := repository.GetText(reminders)
		if len(reminders) == 0 {
			b.Send(m.Sender, "У вас немає нагадувань!")
		} else {
			b.Send(m.Sender, "Ваші нагадування")
			for _, text := range text {
				b.Send(m.Sender, "<b>"+text+"</b>", ModeHTML)
			}
			b.Send(m.Sender, "Введіть номер нагадування, яке потрібно видалити: ")
			b.Handle(tb.OnText, func(m *tb.Message) {
				repository.DeleteEvent(m.Text, b, m)
				b.Send(m.Sender, "Нагадування успішно видалено!")
			})
		}
	})

	b.Handle(&btnDeleteAll, func(m *tb.Message) {
		reminders := repository.GetEvents(m.Sender.ID, services.Timezone.Timezone)

		text := repository.GetText(reminders)
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
					repository.DeleteAllEvents()
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
