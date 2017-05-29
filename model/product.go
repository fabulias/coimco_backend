package model

import "time"

//Product represents the products in the application
type Product struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Details  string `json:"details" binding:"required"`
	Stock    int    `json:"stock" binding:"required"`
	Brand    string `json:"brand" binding:"required"`
	Category string `json:"category" binding:"required"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
