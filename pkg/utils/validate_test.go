package utils

import (
	"testing"

	"github.com/belmadge/freteRapido/pkg/models"
)

func createValidQuoteRequest() models.QuoteRequest {
	return models.QuoteRequest{
		Shipper: models.Shipper{
			RegisteredNumber: "123456",
			Token:            "token",
			PlatformCode:     "platform",
		},
		Recipient: models.Recipient{
			Type:    1,
			Country: "BR",
			Zipcode: 12345,
		},
		Dispatchers: []models.Dispatcher{
			{
				RegisteredNumber: "123456",
				Zipcode:          12345,
				Volumes: []models.Volume{
					{
						Category:      "category",
						Amount:        1,
						UnitaryWeight: 1,
						UnitaryPrice:  1,
						Height:        1,
						Width:         1,
						Length:        1,
					},
				},
			},
		},
	}
}

func TestValidateQuoteInput(t *testing.T) {
	tests := []struct {
		name          string
		input         models.QuoteRequest
		expectedError string
	}{
		{
			name:  "valid input",
			input: createValidQuoteRequest(),
		},
		{
			name: "invalid dispatchers",
			input: func() models.QuoteRequest {
				req := createValidQuoteRequest()
				req.Dispatchers = nil
				return req
			}(),
			expectedError: "recipient zipcode is required",
		},
		{
			name: "no volumes",
			input: func() models.QuoteRequest {
				req := createValidQuoteRequest()
				req.Dispatchers[0].Volumes = []models.Volume{}
				return req
			}(),
			expectedError: "at least one volume is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuoteInput(tt.input)
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error '%v', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error message '%v', got '%v'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got '%v'", err)
				}
			}
		})
	}
}

func TestValidateCarriersFromAPIResponse(t *testing.T) {
	tests := []struct {
		name          string
		apiResponse   map[string]interface{}
		expectedError string
		expectedLen   int
	}{
		{
			name: "Valid API Response",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{
						"offers": []interface{}{
							map[string]interface{}{
								"carrier": map[string]interface{}{
									"name": "Carrier1",
								},
								"final_price": 100.0,
								"service":     "express",
								"delivery_time": map[string]interface{}{
									"days": 3.0,
								},
							},
							map[string]interface{}{
								"carrier": map[string]interface{}{
									"name": "Carrier2",
								},
								"final_price": 200.0,
								"service":     "standard",
								"delivery_time": map[string]interface{}{
									"days": 5.0,
								},
							},
						},
					},
				},
			},
			expectedLen: 2,
		},
		{
			name: "Invalid Dispatchers Data",
			apiResponse: map[string]interface{}{
				"dispatchers": "invalid_data",
			},
			expectedError: "invalid dispatchers data from API response",
		},
		{
			name: "Invalid Dispatcher Format",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					"invalid_format",
				},
			},
			expectedError: "invalid dispatcher format in API response",
		},
		{
			name: "Missing Carrier Name",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{
						"offers": []interface{}{
							map[string]interface{}{
								"final_price": 100.0,
								"service":     "express",
								"delivery_time": map[string]interface{}{
									"days": 3.0,
								},
							},
						},
					},
				},
			},
			expectedError: "missing or invalid dispatcher fields in API response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carriers, err := ValidateCarriersFromAPIResponse(tt.apiResponse)
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error '%v', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error message '%v', got '%v'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got '%v'", err)
				}
				if len(carriers) != tt.expectedLen {
					t.Errorf("expected carriers length %d, got %d", tt.expectedLen, len(carriers))
				}
			}
		})
	}
}
