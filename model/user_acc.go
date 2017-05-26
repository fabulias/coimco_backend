package model

var (
	ADMIN   int8 = 2
	MANAGER int8 = 1
	SELLER  int8 = 0
)

//Represents base of the admin, acquirement manager
//and seller in the application
type User_acc struct {
	Mail     string `json:"mail" db:"name:mail, primarykey" binding:"required"`
	Name     string `json:"name" db:"name:name" binding:"required"`
	Lastname string `json:"lastname" db:"name:lastname" binding:"required"`
	Rut      string `json:"rut" db:"name:rut" binding:"required"`
	Pass     string `json:"pass" db:"name:pass" binding:"required"`
	Role     int8   `json:"role" db:"name:role" binding:"required"`
	Active   bool   `json:"active" db:"name:active"`
}
