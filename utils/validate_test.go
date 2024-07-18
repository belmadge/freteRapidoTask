package utils

import (
	"testing"

	"github.com/belmadge/freteRapido/domain"
	"github.com/stretchr/testify/assert"
)

func TestValidateQuoteInput_Success(t *testing.T) {
	input := domain.QuoteRequest{
		Shipper: domain.Shipper{
			RegisteredNumber: "123456789",
			Token:            "token",
			PlatformCode:     "platform",
		},
		Recipient: domain.Recipient{
			Country: "BRA",
			Zipcode: 12345678,
		},
		Dispatchers: []domain.Dispatcher{
			{
				RegisteredNumber: "123456789",
				Zipcode:          12345678,
				Volumes: []domain.Volume{
					{
						Category:      "7",
						Amount:        1,
						UnitaryWeight: 5,
						UnitaryPrice:  349,
						Height:        0.2,
						Width:         0.2,
						Length:        0.2,
					},
				},
			},
		},
	}

	err := ValidateQuoteInput(input)
	assert.NoError(t, err)
}

func TestValidateQuoteInput_Error(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.QuoteRequest
		wantErr string
	}{
		{
			name: "missing shipper information",
			input: domain.QuoteRequest{
				Shipper: domain.Shipper{},
				Recipient: domain.Recipient{
					Country: "BRA",
					Zipcode: 12345678,
				},
				Dispatchers: []domain.Dispatcher{
					{
						RegisteredNumber: "123456789",
						Zipcode:          12345678,
						Volumes: []domain.Volume{
							{
								Category:      "7",
								Amount:        1,
								UnitaryWeight: 5,
								UnitaryPrice:  349,
								Height:        0.2,
								Width:         0.2,
								Length:        0.2,
							},
						},
					},
				},
			},
			wantErr: "shipper information is incomplete",
		},
		{
			name: "missing recipient information",
			input: domain.QuoteRequest{
				Shipper: domain.Shipper{
					RegisteredNumber: "123456789",
					Token:            "token",
					PlatformCode:     "platform",
				},
				Recipient: domain.Recipient{},
				Dispatchers: []domain.Dispatcher{
					{
						RegisteredNumber: "123456789",
						Zipcode:          12345678,
						Volumes: []domain.Volume{
							{
								Category:      "7",
								Amount:        1,
								UnitaryWeight: 5,
								UnitaryPrice:  349,
								Height:        0.2,
								Width:         0.2,
								Length:        0.2,
							},
						},
					},
				},
			},
			wantErr: "recipient information is incomplete",
		},
		{
			name: "missing dispatchers",
			input: domain.QuoteRequest{
				Shipper: domain.Shipper{
					RegisteredNumber: "123456789",
					Token:            "token",
					PlatformCode:     "platform",
				},
				Recipient: domain.Recipient{
					Country: "BRA",
					Zipcode: 12345678,
				},
			},
			wantErr: "at least one dispatcher is required",
		},
		{
			name: "missing dispatcher information",
			input: domain.QuoteRequest{
				Shipper: domain.Shipper{
					RegisteredNumber: "123456789",
					Token:            "token",
					PlatformCode:     "platform",
				},
				Recipient: domain.Recipient{
					Country: "BRA",
					Zipcode: 12345678,
				},
				Dispatchers: []domain.Dispatcher{
					{
						Volumes: []domain.Volume{
							{
								Category:      "7",
								Amount:        1,
								UnitaryWeight: 5,
								UnitaryPrice:  349,
								Height:        0.2,
								Width:         0.2,
								Length:        0.2,
							},
						},
					},
				},
			},
			wantErr: "dispatcher information is incomplete",
		},
		{
			name: "missing volume information",
			input: domain.QuoteRequest{
				Shipper: domain.Shipper{
					RegisteredNumber: "123456789",
					Token:            "token",
					PlatformCode:     "platform",
				},
				Recipient: domain.Recipient{
					Country: "BRA",
					Zipcode: 12345678,
				},
				Dispatchers: []domain.Dispatcher{
					{
						RegisteredNumber: "123456789",
						Zipcode:          12345678,
					},
				},
			},
			wantErr: "at least one volume is required for each dispatcher",
		},
		{
			name: "invalid volume information",
			input: domain.QuoteRequest{
				Shipper: domain.Shipper{
					RegisteredNumber: "123456789",
					Token:            "token",
					PlatformCode:     "platform",
				},
				Recipient: domain.Recipient{
					Country: "BRA",
					Zipcode: 12345678,
				},
				Dispatchers: []domain.Dispatcher{
					{
						RegisteredNumber: "123456789",
						Zipcode:          12345678,
						Volumes: []domain.Volume{
							{
								Category: "",
								Amount:   0,
							},
						},
					},
				},
			},
			wantErr: "volume information is incomplete or invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuoteInput(tt.input)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateCarriersFromAPIResponse_Success(t *testing.T) {
	tests := []struct {
		name        string
		apiResponse map[string]interface{}
		expected    []domain.Carrier
	}{
		{
			name: "valid response",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{
						"offers": []interface{}{
							map[string]interface{}{
								"carrier": map[string]interface{}{
									"name": "Carrier1",
								},
								"final_price": 10.0,
								"service":     "Service1",
								"delivery_time": map[string]interface{}{
									"days": 2.0,
								},
							},
						},
					},
				},
			},
			expected: []domain.Carrier{
				{Name: "Carrier1", Price: 10.0, Service: "Service1", Deadline: 2},
			},
		},
		{
			name: "valid delivery time in minutes",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{
						"offers": []interface{}{
							map[string]interface{}{
								"carrier": map[string]interface{}{
									"name": "Carrier1",
								},
								"final_price": 10.0,
								"service":     "Service1",
								"delivery_time": map[string]interface{}{
									"minutes": 2880.0,
								},
							},
						},
					},
				},
			},
			expected: []domain.Carrier{
				{Name: "Carrier1", Price: 10.0, Service: "Service1", Deadline: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateCarriersFromAPIResponse(tt.apiResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateCarriersFromAPIResponse_Error(t *testing.T) {
	tests := []struct {
		name        string
		apiResponse map[string]interface{}
		wantErr     string
	}{
		{
			name:        "missing dispatchers",
			apiResponse: map[string]interface{}{},
			wantErr:     "invalid dispatchers data from API response",
		},
		{
			name: "invalid dispatchers list format",
			apiResponse: map[string]interface{}{
				"dispatchers": "invalid",
			},
			wantErr: "invalid dispatchers list format",
		},
		{
			name: "invalid dispatcher format",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					"invalid",
				},
			},
			wantErr: "invalid dispatcher format in API response",
		},
		{
			name: "missing offers in dispatcher",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{},
				},
			},
			wantErr: "missing offers in dispatcher",
		},
		{
			name: "invalid offering format in dispatcher",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{
						"offers": []interface{}{
							"invalid",
						},
					},
				},
			},
			wantErr: "invalid offering format in dispatcher",
		},
		{
			name: "missing carrier in offering",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{
						"offers": []interface{}{
							map[string]interface{}{},
						},
					},
				},
			},
			wantErr: "missing carrier in offering",
		},
		{
			name: "missing or invalid dispatcher fields",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{
						"offers": []interface{}{
							map[string]interface{}{
								"carrier":     map[string]interface{}{"name": "Carrier1"},
								"service":     "Service1",
								"final_price": 10.0,
								"delivery_time": map[string]interface{}{
									"invalid_field": 1,
								},
							},
						},
					},
				},
			},
			wantErr: "missing days, hours, or minutes in delivery_time",
		},
		{
			name: "missing days, hours, or minutes in delivery_time",
			apiResponse: map[string]interface{}{
				"dispatchers": []interface{}{
					map[string]interface{}{
						"offers": []interface{}{
							map[string]interface{}{
								"carrier": map[string]interface{}{
									"name": "Carrier1",
								},
								"final_price":   10.0,
								"service":       "Service1",
								"delivery_time": map[string]interface{}{},
							},
						},
					},
				},
			},
			wantErr: "missing days, hours, or minutes in delivery_time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateCarriersFromAPIResponse(tt.apiResponse)
			assert.EqualError(t, err, tt.wantErr)
			assert.Nil(t, result)
		})
	}
}
