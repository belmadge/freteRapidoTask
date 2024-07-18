package domain

import "time"

type QuoteRequest struct {
	Shipper        Shipper      `json:"shipper"`
	Recipient      Recipient    `json:"recipient"`
	Dispatchers    []Dispatcher `json:"dispatchers"`
	SimulationType []int        `json:"simulation_type"`
}

type Shipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

type Recipient struct {
	Type    int    `json:"type"`
	Country string `json:"country"`
	Zipcode int    `json:"zipcode"`
}

type Dispatcher struct {
	RegisteredNumber string   `json:"registered_number"`
	Zipcode          int      `json:"zipcode"`
	Volumes          []Volume `json:"volumes"`
}

type Volume struct {
	Amount        int     `json:"amount"`
	Category      string  `json:"category"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryPrice  float64 `json:"unitary_price"`
	UnitaryWeight float64 `json:"unitary_weight"`
}

type QuoteResponse struct {
	Carrier []Carrier `json:"carrier"`
}

type Quote struct {
	ID        uint      `gorm:"primaryKey"`
	Carrier   []Carrier `gorm:"foreignKey:QuoteID"`
	CreatedAt time.Time
}

type Carrier struct {
	ID       uint    `gorm:"primaryKey"`
	QuoteID  uint    `gorm:"index"`
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}
