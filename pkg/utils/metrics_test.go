package utils

import (
	"testing"

	"github.com/belmadge/freteRapido/pkg/models"
)

func TestCalculateMetrics(t *testing.T) {
	t.Run("Valid Quotes List", func(t *testing.T) {
		quotes := []models.Quote{
			{
				Carrier: []models.Carrier{
					{Name: "Carrier1", Price: 100},
					{Name: "Carrier2", Price: 200},
				},
			},
			{
				Carrier: []models.Carrier{
					{Name: "Carrier1", Price: 150},
					{Name: "Carrier3", Price: 250},
				},
			},
		}

		metrics, err := CalculateMetrics(quotes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(metrics) == 0 {
			t.Errorf("Expected metrics to be calculated, got empty map")
		}
	})

	t.Run("Empty Quotes List", func(t *testing.T) {
		quotes := []models.Quote{}

		_, err := CalculateMetrics(quotes)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("Quotes with Extreme Values", func(t *testing.T) {
		quotes := []models.Quote{
			{
				Carrier: []models.Carrier{
					{Name: "Carrier1", Price: 1},
					{Name: "Carrier2", Price: 999999},
				},
			},
		}

		metrics, err := CalculateMetrics(quotes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if metrics["cheapest_freight"].(*models.Carrier).Price != 1 {
			t.Errorf("Expected cheapest freight to be 1, got %v", metrics["cheapest_freight"].(*models.Carrier).Price)
		}
		if metrics["most_expensive_freight"].(*models.Carrier).Price != 999999 {
			t.Errorf("Expected most expensive freight to be 999999, got %v", metrics["most_expensive_freight"].(*models.Carrier).Price)
		}
	})
}
