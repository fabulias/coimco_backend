package model

//Represents base of the admin, acquirement manager
//and seller in the application
type User_acc struct {
	Mail     string `json:"mail" db:"mail, primarykey" binding:"required"`
	Name     string `json:"name" db:"name" binding:"required"`
	Lastname string `json:"lastname" db:"lastname" binding:"required"`
	Rut      string `json:"rut" db:"rut" binding:"required"`
	Pass     string `json:"pass" db:"pass" binding:"required"`
	Role     bool   `json:"role" db:"role" binding:"required"`
	Active   bool   `json:"active" db:"active"`
}
