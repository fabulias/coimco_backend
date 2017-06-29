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
	Total uint
}

type ProductPriceID struct {
	Total uint
	Date  time.Time
}

type ProductRankProfitability struct {
	Rent float64
	Name string
}

type ProductRankProviderPrice struct {
	Name  string
	Price uint
	Mail  string
	Phone string
}

type ProductK struct {
	Product
	Cant uint
}
