package utils

import (
	"errors"

	"github.com/belmadge/freteRapido/domain"
)

func CalculateMetrics(quotes []domain.Quote) (map[string]interface{}, error) {
	if len(quotes) == 0 {
		return nil, errors.New("no quotes provided")
	}

	carrierMetrics := make(map[string]map[string]float64)
	var cheapestQuote, mostExpensiveQuote *domain.Carrier

	for _, quote := range quotes {
		for _, carrier := range quote.Carrier {
			updateCarrierMetrics(carrierMetrics, carrier)
			updateCheapestAndMostExpensiveQuote(&cheapestQuote, &mostExpensiveQuote, carrier)
		}
	}

	calculateAveragePrice(carrierMetrics)

	result := map[string]interface{}{
		"carriers":             carrierMetrics,
		"cheapest_quote":       cheapestQuote,
		"most_expensive_quote": mostExpensiveQuote,
	}

	return result, nil
}

func updateCarrierMetrics(carrierMetrics map[string]map[string]float64, carrier domain.Carrier) {
	if _, exists := carrierMetrics[carrier.Name]; !exists {
		carrierMetrics[carrier.Name] = map[string]float64{
			"count":       0,
			"total_price": 0,
		}
	}
	carrierMetrics[carrier.Name]["count"]++
	carrierMetrics[carrier.Name]["total_price"] += carrier.Price
}

func updateCheapestAndMostExpensiveQuote(cheapestQuote, mostExpensiveQuote **domain.Carrier, carrier domain.Carrier) {
	if *cheapestQuote == nil || carrier.Price < (*cheapestQuote).Price {
		*cheapestQuote = &carrier
	}
	if *mostExpensiveQuote == nil || carrier.Price > (*mostExpensiveQuote).Price {
		*mostExpensiveQuote = &carrier
	}
}

func calculateAveragePrice(carrierMetrics map[string]map[string]float64) {
	for _, metrics := range carrierMetrics {
		metrics["average_price"] = metrics["total_price"] / metrics["count"]
	}
}
