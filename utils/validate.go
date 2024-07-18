package utils

import (
	"errors"

	"github.com/belmadge/freteRapido/domain"
)

func ValidateQuoteInput(input domain.QuoteRequest) error {
	if input.Shipper.RegisteredNumber == "" || input.Shipper.Token == "" || input.Shipper.PlatformCode == "" {
		return errors.New("shipper information is incomplete")
	}

	if input.Recipient.Country == "" || input.Recipient.Zipcode == 0 {
		return errors.New("recipient information is incomplete")
	}

	if len(input.Dispatchers) == 0 {
		return errors.New("at least one dispatcher is required")
	}

	for _, dispatcher := range input.Dispatchers {
		if dispatcher.RegisteredNumber == "" || dispatcher.Zipcode == 0 {
			return errors.New("dispatcher information is incomplete")
		}

		if len(dispatcher.Volumes) == 0 {
			return errors.New("at least one volume is required for each dispatcher")
		}

		for _, volume := range dispatcher.Volumes {
			if volume.Category == "" || volume.Amount <= 0 || volume.UnitaryWeight <= 0 || volume.UnitaryPrice <= 0 {
				return errors.New("volume information is incomplete or invalid")
			}
		}
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
			carrier, err := parseCarrier(o)
			if err != nil {
				return nil, err
			}
			carriers = append(carriers, carrier)
		}
	}

	return carriers, nil
}

func parseCarrier(offer interface{}) (domain.Carrier, error) {
	offeringMap, ok := offer.(map[string]interface{})
	if !ok {
		return domain.Carrier{}, errors.New("invalid offering format in dispatcher")
	}

	carrierData, ok := offeringMap["carrier"].(map[string]interface{})
	if !ok {
		return domain.Carrier{}, errors.New("missing carrier in offering")
	}

	carrierName, nameOk := carrierData["name"].(string)
	price, priceOk := offeringMap["final_price"].(float64)
	service, serviceOk := offeringMap["service"].(string)
	deliveryTime, deliveryTimeOk := offeringMap["delivery_time"].(map[string]interface{})

	deadline, err := parseDeliveryTime(deliveryTime)
	if err != nil {
		return domain.Carrier{}, err
	}

	if !nameOk || !priceOk || !serviceOk || !deliveryTimeOk {
		return domain.Carrier{}, errors.New("missing or invalid dispatcher fields in API response")
	}

	return domain.Carrier{
		Name:     carrierName,
		Price:    price,
		Service:  service,
		Deadline: deadline,
	}, nil
}

func parseDeliveryTime(deliveryTime map[string]interface{}) (int, error) {
	if days, ok := deliveryTime["days"].(float64); ok {
		return int(days), nil
	}
	if minutes, ok := deliveryTime["minutes"].(float64); ok {
		return int(minutes / 1440), nil
	}
	if hours, ok := deliveryTime["hours"].(float64); ok {
		return int(hours / 24), nil
	}
	return 0, errors.New("missing days, hours, or minutes in delivery_time")
}
