package domain

type Carrier struct {
	ID       uint    `gorm:"primaryKey"`
	QuoteID  uint    `gorm:"index"`
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}
