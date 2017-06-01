package model

import "time"

type Purchase struct {
	ID         uint      `json:"id"`
	ProviderID string    `json:"id_customer"`
	Date       time.Time `json:"date"`
	ShipTime   time.Time `json:"shiptime"`
}
