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
			name: "invalid shipper registered number",
			input: func() models.QuoteRequest {
				req := createValidQuoteRequest()
				req.Shipper.RegisteredNumber = ""
				return req
			}(),
			expectedError: "shipper registered_number is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuoteInput(tt.input)
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error message '%v', got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
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
						"carrier_name": "Carrier1",
						"price":        100.0,
						"service":      "express",
						"deadline":     3.0,
					},
					map[string]interface{}{
						"carrier_name": "Carrier2",
						"price":        200.0,
						"service":      "standard",
						"deadline":     5.0,
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
						"price":    100.0,
						"service":  "express",
						"deadline": 3.0,
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

func Test_validateDispatcher(t *testing.T) {
	tests := []struct {
		name          string
		input         models.Dispatcher
		expectedError string
	}{
		{
			name: "valid dispatcher",
			input: models.Dispatcher{
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
		{
			name: "invalid dispatcher registered number",
			input: models.Dispatcher{
				RegisteredNumber: "",
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
			expectedError: "dispatcher registered_number is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDispatcher(tt.input)
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error message '%v', got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func Test_validateRecipient(t *testing.T) {
	tests := []struct {
		name          string
		input         models.Recipient
		expectedError string
	}{
		{
			name: "valid recipient",
			input: models.Recipient{
				Type:    1,
				Country: "BR",
				Zipcode: 12345,
			},
		},
		{
			name: "invalid recipient type",
			input: models.Recipient{
				Type:    -1,
				Country: "BR",
				Zipcode: 12345,
			},
			expectedError: "recipient type is required",
		},
		{
			name: "invalid recipient country",
			input: models.Recipient{
				Type:    1,
				Country: "",
				Zipcode: 12345,
			},
			expectedError: "recipient country is required",
		},
		{
			name: "invalid recipient zipcode",
			input: models.Recipient{
				Type:    1,
				Country: "BR",
				Zipcode: 0,
			},
			expectedError: "recipient zipcode is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRecipient(tt.input)
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error message '%v', got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func Test_validateShipper(t *testing.T) {
	tests := []struct {
		name          string
		input         models.Shipper
		expectedError string
	}{
		{
			name: "valid shipper",
			input: models.Shipper{
				RegisteredNumber: "123456",
				Token:            "token",
				PlatformCode:     "platform",
			},
		},
		{
			name: "invalid shipper registered number",
			input: models.Shipper{
				RegisteredNumber: "",
				Token:            "token",
				PlatformCode:     "platform",
			},
			expectedError: "shipper registered_number is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateShipper(tt.input)
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error message '%v', got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func Test_validateVolume(t *testing.T) {
	tests := []struct {
		name          string
		input         models.Volume
		expectedError string
	}{
		{
			name: "valid volume",
			input: models.Volume{
				Category:      "category",
				Amount:        1,
				UnitaryWeight: 1,
				UnitaryPrice:  1,
				Height:        1,
				Width:         1,
				Length:        1,
			},
		},
		{
			name: "invalid volume category",
			input: models.Volume{
				Category:      "",
				Amount:        1,
				UnitaryWeight: 1,
				UnitaryPrice:  1,
				Height:        1,
				Width:         1,
				Length:        1,
			},
			expectedError: "volume category is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVolume(tt.input)
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error message '%v', got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
