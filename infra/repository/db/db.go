package db

import (
	"fmt"

	"github.com/belmadge/freteRapido/config"
	"github.com/belmadge/freteRapido/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	setupDatabaseConnection()
	autoMigrateModels()
}

func setupDatabaseConnection() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.DBUser,
		config.Config.DBPassword,
		config.Config.DBHost,
		config.Config.DBPort,
		config.Config.DBName,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("failed to connect to database:", err)
		return
	}
}

func autoMigrateModels() {
	err := DB.AutoMigrate(&domain.Quote{}, &domain.Carrier{})
	if err != nil {
		logrus.Error("failed to auto-migrate database models:", err)
	}
}
