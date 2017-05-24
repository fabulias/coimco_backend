package model

//Role type
var (
	AdminType       = 0
	AcquirementType = 0
	SellerType      = 1
)

//Represents base of the admin, acquirement manager
//and seller in the application
type User_acc struct {
	Mail     string `json:"mail" binding:"required" db:"mail, primarykey"`
	Name     string `json:"name" binding:"required" db:"name"`
	Lastname string `json:"lastname" binding:"required" db:"lastname"`
	Rut      string `json:"rut" binding:"required" db:"rut"`
	Pass     string `json:"pass" binding:"required" db:"pass"`
	Role     int    `json:"role" binding:"required" db:"role"`
	Active   bool   `db:"active"`
}
