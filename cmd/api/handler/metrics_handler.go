package handler

import (
	"net/http"
	"strconv"

	"github.com/belmadge/freteRapido/domain"
	"github.com/belmadge/freteRapido/infra/repository/db"
	"github.com/belmadge/freteRapido/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const DefaultLastQuotes = 10

// GetMetricsHandler handles the retrieval of metrics based on the quotes stored in the database
func GetMetricsHandler(c *gin.Context) {
	lastQuotesParam := c.Query("last_quotes")

	lastQuotes, err := strconv.Atoi(lastQuotesParam)
	if err != nil || lastQuotes <= 0 {
		logrus.Warn("invalid last_quotes parameter, using default value ", DefaultLastQuotes)
		lastQuotes = DefaultLastQuotes
	}

	var quotes []domain.Quote
	result := db.DB.Preload("Carrier").Order("created_at desc").Limit(lastQuotes).Find(&quotes)
	if result.Error != nil {
		logrus.Error("error fetching quotes: ", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching quotes"})
		return
	}

	metrics, err := utils.CalculateMetrics(quotes)
	if err != nil {
		logrus.Error("failed to calculate metrics: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}