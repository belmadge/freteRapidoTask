package domain

import "time"

type Quote struct {
	ID        uint      `gorm:"primaryKey"`
	Carrier   []Carrier `gorm:"foreignKey:QuoteID"`
	CreatedAt time.Time
}
