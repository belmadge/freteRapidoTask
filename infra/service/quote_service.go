package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/belmadge/freteRapido/domain"
	"github.com/belmadge/freteRapido/utils"
)

// CreateQuote creates a new quote by sending a request to the Frete Rápido API
func CreateQuote(input domain.QuoteRequest) (*domain.QuoteResponse, error) {
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
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://sp.freterapido.com/api/v3/quote/simulate",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get quote from Frete Rápido")
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	var apiResponse map[string]interface{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&apiResponse); err != nil {
		return nil, err
	}

	carriers, err := utils.ValidateCarriersFromAPIResponse(apiResponse)
	if err != nil {
		return nil, err
	}

	quoteResponse := &domain.QuoteResponse{
		Carrier: carriers,
	}

	return quoteResponse, nil
}
