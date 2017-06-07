package model

import "time"

//This struct
type TagCustomer struct {
	TagID      int    `json:"id_tag" gorm:"primary_key"`
	CustomerID string `json:"id_customer" gorm:"primary_key"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
