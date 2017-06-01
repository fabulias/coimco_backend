package model

type PurchaseDetail struct {
	PurchaseID uint `json:"sale_id"`
	ProductID  uint `json:"product_id"`
	Price      uint `json:"price"`
	Quantity   uint `json:"quantity"`
}
