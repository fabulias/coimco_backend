package model

import "time"
import "github.com/jinzhu/gorm"

type Purchase struct {
	gorm.Model
	ProviderID string    `json:"id_provider" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
	ShipTime   time.Time `json:"shiptime" binding:"required"`
}

type PurchaseRankKT struct {
	Name string
	Days uint
}
