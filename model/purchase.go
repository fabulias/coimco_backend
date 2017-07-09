package model

import "time"
import "github.com/jinzhu/gorm"

//This struct represent purchase model
type Purchase struct {
	gorm.Model
	ProviderID string    `json:"id_provider" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
	ShipTime   time.Time `json:"shiptime" binding:"required"`
}

//This struct is to models
type PurchaseRankK struct {
	ProviderName string
	ProductName  string
	Quantity     uint
	Price        uint
	Total        uint
	PurchaseID   uint
}

//This struct is to models
type PurchasesProductRec struct {
	Name  string
	Price uint
	Date  time.Time
}

//This struct is to models
type PurchaseRankProduct struct {
	Cash uint
	Name string
}
