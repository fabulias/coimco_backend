package model

//Product represents the products in the application
type Product struct {
	Id       int    `db:"id, primarykey, autoincrement"`
	Name     string `json:"name" binding:"required" db:"name"`
	Details  string `json:"details" binding:"required" db:"details"`
	Stock    int    `json:"stock" binding:"required" db:"stock"`
	Brand    string `json:"brand" binding:"required" db:"brand"`
	Category string `json:"category" binding:"required" db:"category"`
}
