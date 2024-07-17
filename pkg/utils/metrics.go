package utils

import (
	"errors"

	"github.com/belmadge/freteRapido/pkg/models"
)

func CalculateMetrics(quotes []models.Quote) (map[string]interface{}, error) {
	if quotes == nil || len(quotes) == 0 {
		err := errors.New("no quotes provided")
		return nil, err
	}

	carrierMetrics := make(map[string]map[string]float64)
	var cheapestQuote, mostExpensiveQuote *models.Carrier

	for _, quote := range quotes {
		for _, carrier := range quote.Carrier {
			if _, exists := carrierMetrics[carrier.Name]; !exists {
				carrierMetrics[carrier.Name] = map[string]float64{
					"count":       0,
					"total_price": 0,
				}
			}
			carrierMetrics[carrier.Name]["count"]++
			carrierMetrics[carrier.Name]["total_price"] += carrier.Price

			if cheapestQuote == nil || carrier.Price < cheapestQuote.Price {
				cheapestQuote = &carrier
			}
			if mostExpensiveQuote == nil || carrier.Price > mostExpensiveQuote.Price {
				mostExpensiveQuote = &carrier
			}
		}
	}

	for _, metrics := range carrierMetrics {
		metrics["average_price"] = metrics["total_price"] / metrics["count"]
	}

	result := map[string]interface{}{
		"carriers":             carrierMetrics,
		"cheapest_quote":       cheapestQuote,
		"most_expensive_quote": mostExpensiveQuote,
	}

	return result, nil
}
