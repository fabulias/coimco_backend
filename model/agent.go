package model

import "time"

//Represents base of the clients and providers in the application
type Agent struct {
	Rut   string `json:"rut" binding:"required" gorm:"primary_key;type:varchar(20)"`
	Name  string `json:"name" binding:"required"`
	Mail  string `json:"mail" binding:"required"`
	Phone string `json:"phone"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
