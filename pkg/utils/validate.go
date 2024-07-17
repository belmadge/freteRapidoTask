package utils

import (
	"errors"

	"github.com/belmadge/freteRapido/pkg/models"
)

func ValidateQuoteInput(input models.QuoteRequest) error {
	if err := validateShipper(input.Shipper); err != nil {
		return err
	}

	if err := validateRecipient(input.Recipient); err != nil {
		return err
	}

	if err := validateDispatchers(input.Dispatchers); err != nil {
		return err
	}

	return nil
}

func ValidateCarriersFromAPIResponse(apiResponse map[string]interface{}) ([]models.Carrier, error) {
	dispatchersData, ok := apiResponse["dispatchers"].([]interface{})
	if !ok {
		return nil, errors.New("invalid dispatchers data from API response")
	}

	var carriers []models.Carrier
	for _, d := range dispatchersData {
		dispatcherMap, ok := d.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid dispatcher format in API response")
		}

		name, nameOk := dispatcherMap["carrier_name"].(string)
		price, priceOk := dispatcherMap["price"].(float64)
		service, serviceOk := dispatcherMap["service"].(string)
		deadlineFloat, deadlineOk := dispatcherMap["deadline"].(float64)

		if !nameOk || !priceOk || !serviceOk || !deadlineOk {
			return nil, errors.New("missing or invalid dispatcher fields in API response")
		}

		carrier := models.Carrier{
			Name:     name,
			Price:    price,
			Service:  service,
			Deadline: int(deadlineFloat),
		}

		carriers = append(carriers, carrier)
	}

	return carriers, nil
}

func validateShipper(shipper models.Shipper) error {
	if shipper.RegisteredNumber == "" {
		return errors.New("shipper registered_number is required")
	}

	if shipper.Token == "" {
		return errors.New("shipper token is required")
	}

	if shipper.PlatformCode == "" {
		return errors.New("shipper platform_code is required")
	}

	return nil
}

func validateRecipient(recipient models.Recipient) error {
	if recipient.Type < 0 || recipient.Type > 1 {
		return errors.New("recipient type is required")
	}

	if recipient.Country == "" {
		return errors.New("recipient country is required")
	}

	if recipient.Zipcode == 0 {
		return errors.New("recipient zipcode is required")
	}

	return nil
}

func validateDispatchers(dispatchers []models.Dispatcher) error {
	if len(dispatchers) == 0 {
		return errors.New("at least one dispatcher is required")
	}

	for _, dispatcher := range dispatchers {
		if err := validateDispatcher(dispatcher); err != nil {
			return err
		}
	}

	return nil
}

func validateDispatcher(dispatcher models.Dispatcher) error {
	if dispatcher.RegisteredNumber == "" {
		return errors.New("dispatcher registered_number is required")
	}

	if dispatcher.Zipcode == 0 {
		return errors.New("dispatcher zipcode is required")
	}

	if len(dispatcher.Volumes) == 0 {
		return errors.New("at least one volume is required for each dispatcher")
	}

	for _, volume := range dispatcher.Volumes {
		if err := validateVolume(volume); err != nil {
			return err
		}
	}

	return nil
}

func validateVolume(volume models.Volume) error {
	if volume.Category == "" {
		return errors.New("volume category is required")
	}

	if volume.Amount <= 0 {
		return errors.New("volume amount is required")
	}

	if volume.UnitaryWeight <= 0 {
		return errors.New("volume unitary weight is required")
	}

	if volume.UnitaryPrice <= 0 {
		return errors.New("volume unitary price is required")
	}

	if volume.Height <= 0 {
		return errors.New("volume height is required")
	}

	if volume.Width <= 0 {
		return errors.New("volume width is required")
	}

	if volume.Length <= 0 {
		return errors.New("volume length is required")
	}

	return nil
}
