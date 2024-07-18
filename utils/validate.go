package utils

import (
	"errors"

	"github.com/belmadge/freteRapido/domain"
)

func ValidateQuoteInput(input domain.QuoteRequest) error {
	if input.Dispatchers == nil {
		return errors.New("recipient zipcode is required")
	}

	if len(input.Dispatchers) == 0 {
		return errors.New("at least one volume is required")
	}

	return nil
}

func ValidateCarriersFromAPIResponse(apiResponse map[string]interface{}) ([]domain.Carrier, error) {
	dispatchersData, ok := apiResponse["dispatchers"]
	if !ok {
		return nil, errors.New("invalid dispatchers data from API response")
	}

	dispatchersList, ok := dispatchersData.([]interface{})
	if !ok {
		return nil, errors.New("invalid dispatchers list format")
	}

	var carriers []domain.Carrier
	for _, d := range dispatchersList {
		dispatcherMap, ok := d.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid dispatcher format in API response")
		}

		offers, ok := dispatcherMap["offers"].([]interface{})
		if !ok {
			return nil, errors.New("missing offers in dispatcher")
		}

		for _, o := range offers {
			offeringMap, ok := o.(map[string]interface{})
			if !ok {
				return nil, errors.New("invalid offering format in dispatcher")
			}

			carrierData, ok := offeringMap["carrier"].(map[string]interface{})
			if !ok {
				return nil, errors.New("missing carrier in offering")
			}

			carrierName, nameOk := carrierData["name"].(string)
			price, priceOk := offeringMap["final_price"].(float64)
			service, serviceOk := offeringMap["service"].(string)
			deliveryTime, deliveryTimeOk := offeringMap["delivery_time"].(map[string]interface{})

			var deadline int
			var deadlineOk bool

			if days, daysOk := deliveryTime["days"].(float64); daysOk {
				deadline = int(days)
				deadlineOk = true
			} else if minutes, minutesOk := deliveryTime["minutes"].(float64); minutesOk {
				deadline = int(minutes / 1440) // Convertendo minutos para dias
				deadlineOk = true
			} else if hours, hoursOk := deliveryTime["hours"].(float64); hoursOk {
				deadline = int(hours / 24) // Convertendo horas para dias
				deadlineOk = true
			} else {
				return nil, errors.New("missing days, hours, or minutes in delivery_time")
			}

			if !nameOk || !priceOk || !serviceOk || !deliveryTimeOk || !deadlineOk {
				return nil, errors.New("missing or invalid dispatcher fields in API response")
			}

			carrier := domain.Carrier{
				Name:     carrierName,
				Price:    price,
				Service:  service,
				Deadline: deadline,
			}

			carriers = append(carriers, carrier)
		}
	}

	return carriers, nil
}
