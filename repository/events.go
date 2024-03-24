package repository

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
	"log"
	"remider/database"
	"remider/models"
	"strconv"
	"time"
)

const (
	ModeHTML tb.ParseMode = "HTML"
)

func CreateEvent(event *models.Event) error {
	if err := database.DB.Create(event).Error; err != nil {
		fmt.Println("error: ", err)
		return err
	}
	return nil
}

func CheckEvents(b *tb.Bot, timezone string) {

	logrus.Println("проверка напоминаний...")

	var reminders []models.Event

	currentTime := time.Now().UTC()

	database.DB.Where("date <= ?", currentTime).Find(&reminders)

	for _, reminder := range reminders {

		userLocation, err := time.LoadLocation(timezone)
		if err != nil {
			logrus.Fatal(err)
		}

		userTime := reminder.Date.In(userLocation)

		if userTime.Before(currentTime) || userTime.Equal(currentTime) {
			_, err := b.Send(&tb.Chat{ID: reminder.UserId}, "Нагадування: <b>"+reminder.Text+"</b>", ModeHTML)
			if err != nil {
				log.Printf("Не удалось отправить сообщение пользователю: %v", err)
				continue
			}

			if err := database.DB.Delete(&reminder).Error; err != nil {
				log.Printf("Не удалось удалить напоминание из базы данных: %v", err)
			} else {
				log.Printf("Напоминание '%s' удалено из базы данных", reminder.Text)
			}
		}
	}
}

func GetEvents(userid int64, timezone string) []models.Event {
	if timezone == "" {
		return []models.Event{}
	}

	var reminders []models.Event
	database.DB.Where("user_id = ?", userid).Find(&reminders)

	for i, reminder := range reminders {
		userLocation, err := time.LoadLocation(timezone)
		if err != nil {
			logrus.Printf("Ошибка загрузки часового пояса пользователя: %v\n", err)
			continue
		}

		reminders[i].Date = reminder.Date.In(userLocation)
	}

	return reminders
}

func UpdateEvent(reminder *models.Event) error {
	if err := database.DB.Save(reminder).Error; err != nil {
		return err
	}
	return nil
}

func GetText(reminders []models.Event) []string {
	var texts []string

	for _, reminder := range reminders {
		userLocation, err := time.LoadLocation(reminder.Timezone)
		if err != nil {
			logrus.Fatal(err)
		}

		reminderTime := reminder.Date.In(userLocation)

		textWithDate := fmt.Sprintf("%d. %s - %s", reminder.ID, reminder.Text, reminderTime.Format("2006-01-02 15:04")+"\n"+reminder.Timezone)

		texts = append(texts, textWithDate)
	}

	return texts
}

func DeleteEvent(numberString string, b *tb.Bot, m *tb.Message) error {
	var reminder models.Event

	number, err := strconv.Atoi(numberString)
	if err != nil {
		return err
	}

	if err = database.DB.First(&reminder, number).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			b.Send(m.Sender, "reminder with id %d not found", number)
			return fmt.Errorf("reminder with id %d not found", number)
		}
		return err
	}

	database.DB.Unscoped().Delete(&reminder)

	return nil
}

func DeleteAllEvents() error {
	var reminders []models.Event

	if err := database.DB.Migrator().DropTable(&reminders); err != nil {
		return err
	}

	if err := database.DB.Migrator().CreateTable(&reminders); err != nil {
		return err
	}
	return nil
}
