package model

import "github.com/jinzhu/gorm"
import "time"

//Product represents the products in the application
type Product struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Details  string `json:"details" binding:"required"`
	Brand    string `json:"brand" binding:"required"`
	Category string `json:"category" binding:"required"`
}

type InfoProduct struct {
	ID    uint
	Name  string
	Sales uint
	Total uint
}

type ProductRankCategory struct {
	ID    uint
	Name  string
	Sales uint
}

type ProductPrice struct {
	Price uint
	Date  time.Time
}
