package model

import "time"
import "github.com/jinzhu/gorm"

type Sale struct {
	gorm.Model
	CustomerID string    `json:"id_customer" binding:"required"`
	UserID     string    `json:"id_user" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
}

type InfoDashboard struct {
	Count uint
	Sum   uint
}

type SaleRankK struct {
	ID   uint
	Name string
	Cash uint
}

type SaleRankCategory struct {
	ID   uint
	Name string
	Cash uint
}

type SaleRankProduct struct {
	Cash uint
	Name string
}

type SaleRankArea struct {
	Name string
	Cash uint
}

type SaleProductPrice struct {
	Date  time.Time
	Price uint
}
