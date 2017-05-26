package model

//Product represents the products in the application
type Product struct {
	Id       int    `db:"name:id, primarykey, autoincrement"`
	Name     string `json:"name" binding:"required" db:"name:name"`
	Details  string `json:"details" binding:"required" db:"name:details"`
	Stock    int    `json:"stock" binding:"required" db:"name:stock"`
	Brand    string `json:"brand" binding:"required" db:"name:brand"`
	Category string `json:"category" binding:"required" db:"name:category"`
}
