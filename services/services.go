package services

import (
	"encoding/json"
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"remider/config"
	"remider/models"
	"remider/repository"
	"strconv"
	"time"
)

var (
	text      string
	userId    int64
	firstname string
	Timezone  models.TimeZoneInfo
)

func SaveEvent(b *tb.Bot, m *tb.Message) {
	if Timezone.Timezone == "" {
		b.Send(m.Sender, "Ви не вписали свое мiсто для вибору часового поясу!\n"+
			"Щоб вписати свое мiсто нажмiть на команду /setcity")
		return
	}

	userLocation := GetLocation(Timezone.Timezone)
	date, err := time.ParseInLocation("2006-01-02 15:04", m.Text, userLocation)
	if err != nil {
		b.Send(m.Sender, "Неправильний формат дати, спробуй ще раз!")
		return
	}
	if date.Before(time.Now()) {
		b.Send(m.Sender, "Дата не може бути раніше за поточну, спробуй ще раз!")
		return
	}

	userDate := date.UTC()

	event := models.Event{
		Firstname: firstname,
		Date:      userDate,
		UserId:    userId,
		Text:      text,
		Timezone:  Timezone.Timezone,
	}

	if err := repository.CreateEvent(&event); err != nil {
		logrus.Fatal(err)
	}
	b.Send(m.Sender, "Подія успішно додана до бази даних! :)")
}

func GetLocation(timezone string) *time.Location {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logrus.Println("Ошибка загрузки часового пояса:", err)
		return time.UTC
	}
	return loc
}

func SaveEventText(m *tb.Message) {
	text = m.Text
	firstname = m.Sender.FirstName
	userId = m.Sender.ID
}

func EditTimeForEvent(b *tb.Bot, m *tb.Message) {
	reminders := repository.GetEvents(m.Sender.ID, Timezone.Timezone)

	reminderID, err := strconv.Atoi(m.Text)
	if err != nil {
		b.Send(m.Sender, "Неправильний формат нагадування!")
		return
	}

	foundReminder := findReminderByID(reminders, reminderID)

	if foundReminder != nil {
		b.Send(m.Sender, "Введіть нову дату нагадування 2023-12-31 23:59: ")

		b.Handle(tb.OnText, func(m *tb.Message) {
			updateReminderDate(b, m, foundReminder)
		})
	} else {
		b.Send(m.Sender, "У вас немає нагадування із таким номером!")
	}
}

func EditTextForEvent(b *tb.Bot, m *tb.Message) {
	reminders := repository.GetEvents(m.Sender.ID, Timezone.Timezone)

	reminderID, err := strconv.Atoi(m.Text)
	if err != nil {
		b.Send(m.Sender, "Неправильний формат нагадування!")
		return
	}

	foundReminder := findReminderByID(reminders, reminderID)

	if foundReminder != nil {
		b.Send(m.Sender, "Введіть новий текст для нагадування: ")

		b.Handle(tb.OnText, func(m *tb.Message) {
			newText := m.Text

			foundReminder.Text = newText

			if err = repository.UpdateEvent(foundReminder); err != nil {
				log.Fatal(err)
			}

			b.Send(m.Sender, "Нагадування успішно змінено!")
		})
	} else {
		b.Send(m.Sender, "У вас немає нагадування із таким номером!")
	}

}

func GetCity(city string) {
	result, _ := gt.Translate(city, "auto", "en")
	if result == "NY" {
		result = "NewYork"
	}

	url := fmt.Sprintf(config.Cfg.TzApiUrl, result)

	res, err := http.Get(url)
	if err != nil {
		logrus.Fatal(err)
	}

	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&Timezone); err != nil {
		logrus.Fatal(err)
	}
}

func findReminderByID(reminders []models.Event, reminderID int) *models.Event {
	for i := range reminders {
		if int(reminders[i].ID) == reminderID {
			return &reminders[i]
		}
	}
	return nil
}

func updateReminderDate(b *tb.Bot, m *tb.Message, foundReminder *models.Event) {
	if Timezone.Timezone == "" {
		b.Send(m.Sender, "Ви не вписали свое мiсто для вибору часового поясу!\n"+
			"Щоб вписати свое мiсто нажмiть на команду /setcity")
		return
	}

	userLocation := GetLocation(Timezone.Timezone)

	date, err := time.ParseInLocation("2006-01-02 15:04", m.Text, userLocation)
	if err != nil {
		b.Send(m.Sender, "Неправильний формат дати, спробуй ще раз!")
		return
	}
	if date.Before(time.Now()) {
		b.Send(m.Sender, "Дата не може бути раніше за поточну, спробуй ще раз!")
		return
	}

	userDate := date.UTC()

	foundReminder.Date = userDate

	err = repository.UpdateEvent(foundReminder)
	if err != nil {
		log.Fatal(err)
	}

	b.Send(m.Sender, "Дата нагадування успішно змінено!")
}
