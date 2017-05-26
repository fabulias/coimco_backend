package model

//Represents base of the clients and providers in the application
type Agent struct {
	Rut   string `json:"rut" db:"name: rut, primarykey" binding:"required"`
	Name  string `json:"name" db:"name:name" binding:"required"`
	Mail  string `json:"mail" db:"name:mail" binding:"required"`
	Phone string `json:"phone" db:"name:phone"`
}
