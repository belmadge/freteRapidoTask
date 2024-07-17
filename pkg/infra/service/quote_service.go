package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/belmadge/freteRapido/pkg/models"
	"github.com/belmadge/freteRapido/pkg/utils"
	"github.com/sirupsen/logrus"
)

// CreateQuote creates a new quote by sending a request to the Frete Rápido API
func CreateQuote(input models.QuoteRequest) (*models.QuoteResponse, error) {
	if err := utils.ValidateQuoteInput(input); err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"shipper": map[string]string{
			"registered_number": input.Shipper.RegisteredNumber,
			"token":             input.Shipper.Token,
			"platform_code":     input.Shipper.PlatformCode,
		},
		"recipient": map[string]interface{}{
			"type":    input.Recipient.Type,
			"country": input.Recipient.Country,
			"zipcode": input.Recipient.Zipcode,
		},
		"dispatchers":     input.Dispatchers,
		"simulation_type": input.SimulationType,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		logrus.Error("failed to create payload:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://sp.freterapido.com/api/v3/quote/simulate",
		bytes.NewBuffer(requestBody))
	if err != nil {
		logrus.Error("failed to create request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("failed to send request:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error("failed to get quote from Frete Rápido")
		return nil, errors.New("failed to get quote from Frete Rápido")
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	var apiResponse map[string]interface{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&apiResponse); err != nil {
		logrus.Error("failed to decode response:", err)
		return nil, err
	}

	carriers, err := utils.ValidateCarriersFromAPIResponse(apiResponse)
	if err != nil {
		logrus.Error("failed to validate carriers:", err)
		return nil, err
	}

	quoteResponse := &models.QuoteResponse{
		Carrier: carriers,
	}

	return quoteResponse, nil
}
