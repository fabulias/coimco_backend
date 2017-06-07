package model

import "github.com/jinzhu/gorm"

//Product represents the products in the application
type Product struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Details  string `json:"details" binding:"required"`
	Brand    string `json:"brand" binding:"required"`
	Category string `json:"category" binding:"required"`
}
