package model

type Agent struct {
	Rut   string `json:"rut" db:"rut, primarykey" binding:"required"`
	Name  string `json:"name" db:"name" binding:"required"`
	Mail  string `json:"mail" db:"mail" binding:"required"`
	Phone string `json:"phone" db:"phone"`
}
