package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}

	Config.DBUser = os.Getenv("DB_USER")
	Config.DBPassword = os.Getenv("DB_PASSWORD")
	Config.DBName = os.Getenv("DB_NAME")
	Config.DBHost = os.Getenv("DB_HOST")
	Config.DBPort = os.Getenv("DB_PORT")
}
