package main

import (
	"github.com/belmadge/freteRapido/config"
	"github.com/belmadge/freteRapido/pkg/handler"
	"github.com/belmadge/freteRapido/pkg/infra/repository/db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()
	db.InitDB()

	r := gin.Default()

	r.POST("/quote", handler.CreateQuoteHandler)
	r.GET("/metrics", handler.GetMetricsHandler)

	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %s", err.Error())
	}
}
