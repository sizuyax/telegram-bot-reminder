package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"remider/models"
)

var Cfg *models.Config

func LoadEnv() {

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("failed to load .env file")
	}

	Cfg = &models.Config{}

	if err := env.Parse(Cfg); err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugln("successfully loaded .env file")
}
