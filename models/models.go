package models

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	Firstname string    `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
	Text      string    `gorm:"not null"`
	UserId    int64     `gorm:"not null"`
	Timezone  string    `gorm:"not null"`
}

type TimeZoneInfo struct {
	Geo               GeoInfo `json:"geo"`
	Timezone          string  `json:"timezone"`
	DateTime          string  `json:"date_time"`
	TimezoneOffset    int     `json:"timezone_offset"`
	TimezoneOffsetDST int     `json:"timezone_offset_with_dst"`
	Date              string  `json:"date"`
}

type GeoInfo struct {
	Location  string  `json:"location"`
	Country   string  `json:"country"`
	State     string  `json:"state"`
	City      string  `json:"city"`
	Locality  string  `json:"locality"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Config struct {
	TgBotToken       string `env:"TG_BOT_TOKEN"`
	TzApiUrl         string `env:"TZ_API_URL"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDBName   string `env:"POSTGRES_DB"`
	PostgresPorts    string `env:"POSTGRES_PORTS"`
	PostgresHost     string `env:"POSTGRES_HOST"`
}
