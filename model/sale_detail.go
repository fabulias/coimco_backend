package model

type SaleDetail struct {
	SaleID    uint `json:"sale_id" binding:"required" gorm:"primary_key"`
	ProductID uint `json:"product_id" binding:"required" gorm:"primary_key"`
	Price     uint `json:"price" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required"`
}
