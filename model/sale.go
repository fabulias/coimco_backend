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
