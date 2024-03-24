package database

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"remider/config"
	"remider/models"
)

var DB *gorm.DB

func InitDB() error {
	var err error

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.Cfg.PostgresUser,
		config.Cfg.PostgresPassword,
		config.Cfg.PostgresHost,
		config.Cfg.PostgresDBName,
		"disable",
	)
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		logrus.Error("error: ", err)
		return err
	}

	fmt.Println("db connected")

	if err = DB.AutoMigrate(&models.Event{}); err != nil {
		fmt.Println("error: ", err)
		return err
	}

	return nil
}
