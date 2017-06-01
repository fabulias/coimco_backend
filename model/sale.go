package model

import "time"

type Sale struct {
	ID         uint      `json:"id"`
	CustomerID string    `json:"id_customer"`
	UserID     string    `json:"id_user"`
	Date       time.Time `json:"date"`
}
