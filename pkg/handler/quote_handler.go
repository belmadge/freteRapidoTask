package handler

import (
	"net/http"

	"github.com/belmadge/freteRapido/pkg/infra/repository/db"
	"github.com/belmadge/freteRapido/pkg/infra/service"
	"github.com/belmadge/freteRapido/pkg/models"
	"github.com/belmadge/freteRapido/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreateQuoteHandler handles the creation of a new quote
func CreateQuoteHandler(c *gin.Context) {
	var input models.QuoteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error("error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateQuoteInput(input); err != nil {
		logrus.Error("error validating quote input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quoteResponse, err := service.CreateQuote(input)
	if err != nil {
		logrus.Error("error creating quote: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quote := models.Quote{
		Carrier: quoteResponse.Carrier,
	}

	result := db.DB.Create(&quote)
	if result.Error != nil {
		logrus.Error("error saving quote to database: ", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving quote to database"})
		return
	}

	c.JSON(http.StatusCreated, quoteResponse)
}
