package model

import "time"

var (
	ADMIN   int8 = 2
	MANAGER int8 = 1
	SELLER  int8 = 0
)

//Represents base of the admin, acquirement manager
//and seller in the application
type User_acc struct {
	Mail     string `json:"mail" binding:"required" gorm:"primary_key"`
	Name     string `json:"name" binding:"required"`
	Lastname string `json:"lastname" binding:"required"`
	Rut      string `json:"rut" binding:"required"`
	Pass     string `json:"pass" binding:"required"`
	Role     int8   `json:"role" binding:"required"`
	Active   bool   `json:"active" `

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
