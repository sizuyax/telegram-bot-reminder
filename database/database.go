package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

type Event struct {
	gorm.Model
	Firstname string    `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
	Text      string    `gorm:"not null"`
	User_id   int64     `gorm:"not null"`
}

func InitDB() error {
	dsn := "host=localhost user=admin password=pass dbname=db_auth port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}

	fmt.Println("db connected", db)

	err = db.AutoMigrate(&Event{})
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}

	return nil
}

func CreateEvent(event *Event) error {
	err := db.Create(event).Error
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}
	return nil
}
