package model

import "time"

//This struct
type TagCustomer struct {
	TagID      int    `json:"id_tag"`
	CustomerID string `json:"id_customer"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
