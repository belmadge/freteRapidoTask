package utils

import (
	"testing"

	"github.com/belmadge/freteRapido/domain"
	"github.com/stretchr/testify/assert"
)

func TestCalculateMetrics_Success(t *testing.T) {
	quotes := []domain.Quote{
		{
			Carrier: []domain.Carrier{
				{Name: "Carrier1", Price: 10},
				{Name: "Carrier2", Price: 20},
			},
		},
		{
			Carrier: []domain.Carrier{
				{Name: "Carrier1", Price: 30},
				{Name: "Carrier2", Price: 40},
			},
		},
	}

	expected := map[string]interface{}{
		"carriers": map[string]map[string]float64{
			"Carrier1": {"count": 2, "total_price": 40, "average_price": 20},
			"Carrier2": {"count": 2, "total_price": 60, "average_price": 30},
		},
		"cheapest_quote":       &domain.Carrier{Name: "Carrier1", Price: 10},
		"most_expensive_quote": &domain.Carrier{Name: "Carrier2", Price: 40},
	}

	result, err := CalculateMetrics(quotes)

	assert.Nil(t, err)
	assert.Equal(t, expected["carriers"], result["carriers"])
	assert.Equal(t, expected["cheapest_quote"].(*domain.Carrier).Name, result["cheapest_quote"].(*domain.Carrier).Name)
	assert.Equal(t, expected["cheapest_quote"].(*domain.Carrier).Price, result["cheapest_quote"].(*domain.Carrier).Price)
	assert.Equal(t, expected["most_expensive_quote"].(*domain.Carrier).Name, result["most_expensive_quote"].(*domain.Carrier).Name)
	assert.Equal(t, expected["most_expensive_quote"].(*domain.Carrier).Price, result["most_expensive_quote"].(*domain.Carrier).Price)
}

func TestCalculateMetrics_Error(t *testing.T) {
	quotes := []domain.Quote{}

	_, err := CalculateMetrics(quotes)

	assert.NotNil(t, err)
	assert.Equal(t, "no quotes provided", err.Error())
}
