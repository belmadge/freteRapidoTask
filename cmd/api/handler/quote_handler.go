package handler

import (
	"net/http"

	"github.com/belmadge/freteRapido/domain"
	"github.com/belmadge/freteRapido/infra/repository/db"
	"github.com/belmadge/freteRapido/infra/service"
	"github.com/belmadge/freteRapido/utils"
	"github.com/gin-gonic/gin"
)

// CreateQuoteHandler handles the creation of a new quote
func CreateQuoteHandler(c *gin.Context) {
	var input domain.QuoteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateQuoteInput(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quoteResponse, err := service.CreateQuote(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quote := domain.Quote{
		Carrier: quoteResponse.Carrier,
	}

	result := db.DB.Create(&quote)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving quote to database"})
		return
	}

	c.JSON(http.StatusCreated, quoteResponse)
}
