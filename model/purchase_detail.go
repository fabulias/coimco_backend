package model

type PurchaseDetail struct {
	PurchaseID uint `json:"purchase_id" binding:"required"`
	ProductID  uint `json:"product_id" binding:"required"`
	Price      uint `json:"price" binding:"required"`
	Quantity   uint `json:"quantity" binding:"required"`
}
