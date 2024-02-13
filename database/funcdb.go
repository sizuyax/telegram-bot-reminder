package database

import (
	"errors"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

func CreateEvent(event *Event, m *tb.Message) error {
	event.Event_id = m.Sender.ID

	err := db.Create(event).Error
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}
	return nil
}

func CheckReminders(b *tb.Bot) {

	log.Println("проверка напоминаний...")
	var reminders []Event
	currentTime := time.Now().UTC()

	db.Where("date <= ?", currentTime).Find(&reminders)
	for _, reminder := range reminders {
		_, err := b.Send(&tb.Chat{ID: reminder.User_id}, "Нагадування: <b>"+reminder.Text+"</b>", ModeHTML)
		if err != nil {
			log.Printf("Failed to send message to user: %v", err, reminder.User_id)
			continue
		}

		log.Printf("Deleting reminder '%s' from database", reminder.Text)
		if err := db.Delete(&reminder).Error; err != nil {
			log.Printf("Failed to delete reminder from database: %v", err)
		}
	}
}

func GetReminders(userid int64) []Event {
	var reminders []Event
	db.Where("user_id = ?", userid).Find(&reminders)
	return reminders
}

func UpdateReminder(reminder *Event) error {
	err := db.Save(reminder).Error
	if err != nil {
		return err
	}
	return nil
}

func GetText() []string {
	//var reminder Event
	var texts []string
	var reminders []Event

	if err := db.Model(&reminders).Find(&reminders).Error; err != nil {
		// Если произошла ошибка при получении событий, вернуть пустой срез строк
		return texts
	}

	// Преобразовать полученные события в тексты с датами
	for _, reminder := range reminders {
		textWithDate := fmt.Sprintf("%d. %s - %s", reminder.ID, reminder.Text, reminder.Date)
		texts = append(texts, textWithDate)
	}

	// Вернуть все тексты с датами
	return texts
}

func DeleteReminder(numberString string, b *tb.Bot, m *tb.Message) error {
	var reminder Event

	number, err := strconv.Atoi(numberString)
	if err != nil {
		return err
	}

	err = db.First(&reminder, number).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			b.Send(m.Sender, "reminder with id %d not found", number)
			return fmt.Errorf("reminder with id %d not found", number)
		}
		return err
	}

	db.Unscoped().Delete(&reminder)

	return nil
}

func DeleteAllReminders() error {
	var reminders []Event

	err := db.Migrator().DropTable(&reminders)
	if err != nil {
		return err
	}

	err = db.Migrator().CreateTable(&reminders)
	if err != nil {
		return err
	}
	return nil
}
