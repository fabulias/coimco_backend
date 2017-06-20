package model

type Provider struct {
	Agent
}

type ProviderRankK struct {
	ProviderName string "gorm:providerName"
	ProductName  string "gorm:productName"
	Quantity     uint
	Price        uint
	Total        uint
	PurchaseID   uint "gorm:id_purchase"
}
