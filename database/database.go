package database

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

// Event struct
type Event struct {
	gorm.Model
	Firstname string    `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
	Text      string    `gorm:"not null"`
	User_id   int64     `gorm:"not null"`
	Event_id  int64     `gorm:"not null"`
}

const (
	ModeHTML tb.ParseMode = "HTML"
)

// func for initializing db
func InitDB() (*gorm.DB, error) {
	dsn := "postgres://qxiukmkb:zDC-9sAEJOU6IFQB5ga5G6Bn6fAwNYrV@dumbo.db.elephantsql.com/qxiukmkb"
	var err error
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	fmt.Println("db connected", db)

	err = db.AutoMigrate(&Event{})
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	return db, nil
}
