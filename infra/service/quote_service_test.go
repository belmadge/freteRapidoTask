package service

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/belmadge/freteRapido/domain"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateQuote_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://sp.freterapido.com/api/v3/quote/simulate",
		httpmock.NewStringResponder(200, `{
			"dispatchers": [{
				"offers": [{
					"carrier": {"name": "Carrier1"},
					"final_price": 10.0,
					"service": "Service1",
					"delivery_time": {"days": 2}
				}]
			}]
		}`))

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
		SimulationType: []int{0},
	}

	expectedResponse := domain.QuoteResponse{
		Carrier: []domain.Carrier{
			{Name: "Carrier1", Price: 10.0, Service: "Service1", Deadline: 2},
		},
	}

	quoteResponse, err := CreateQuote(input)

	assert.NoError(t, err)
	assert.Equal(t, &expectedResponse, quoteResponse)
}

func TestCreateQuote_Error(t *testing.T) {
	tests := []struct {
		name        string
		input       domain.QuoteRequest
		setupMock   func()
		expectedErr string
	}{
		{
			name: "validation error",
			input: domain.QuoteRequest{
				Shipper: domain.Shipper{},
			},
			setupMock:   func() {},
			expectedErr: "shipper information is incomplete",
		},
		{
			name: "http request error",
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
				SimulationType: []int{0},
			},
			setupMock: func() {
				httpmock.Activate()
				httpmock.RegisterResponder("POST", "https://sp.freterapido.com/api/v3/quote/simulate",
					httpmock.NewErrorResponder(errors.New("http request error")))
			},
			expectedErr: "Post \"https://sp.freterapido.com/api/v3/quote/simulate\": http request error",
		},
		{
			name: "non-200 response",
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
				SimulationType: []int{0},
			},
			setupMock: func() {
				httpmock.Activate()
				httpmock.RegisterResponder("POST", "https://sp.freterapido.com/api/v3/quote/simulate",
					httpmock.NewStringResponder(500, `{}`))
			},
			expectedErr: "failed to get quote from Frete Rápido",
		},
		{
			name: "invalid response body",
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
				SimulationType: []int{0},
			},
			setupMock: func() {
				httpmock.Activate()
				httpmock.RegisterResponder("POST", "https://sp.freterapido.com/api/v3/quote/simulate",
					httpmock.NewStringResponder(200, `invalid`))
			},
			expectedErr: "invalid character 'i' looking for beginning of value",
		},
		{
			name: "validation response error",
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
				SimulationType: []int{0},
			},
			setupMock: func() {
				httpmock.Activate()
				apiResponse := map[string]interface{}{
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
										"invalid_field": 1,
									},
								},
							},
						},
					},
				}
				responseBody, _ := json.Marshal(apiResponse)
				httpmock.RegisterResponder("POST", "https://sp.freterapido.com/api/v3/quote/simulate",
					httpmock.NewStringResponder(200, string(responseBody)))
			},
			expectedErr: "missing days, hours, or minutes in delivery_time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			defer httpmock.DeactivateAndReset()

			quoteResponse, err := CreateQuote(tt.input)

			assert.Nil(t, quoteResponse)
			assert.EqualError(t, err, tt.expectedErr)
			httpmock.DeactivateAndReset()
		})
	}
}
