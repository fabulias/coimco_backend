package model

import "time"
import "github.com/jinzhu/gorm"

type Purchase struct {
	gorm.Model
	ProviderID string    `json:"id_provider" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
	ShipTime   time.Time `json:"shiptime" binding:"required"`
}

type PurchaseRankK struct {
	ProviderName string //"gorm:providerName"
	ProductName  string //"gorm:productName"
	Quantity     uint
	Price        uint
	Total        uint
	PurchaseID   uint //"gorm:id_purchase"
}
