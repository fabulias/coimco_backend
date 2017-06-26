package model

//GetPurchasesProduct returns history of a product in purchases
func GetPurchasesProduct(id string, in Date) ([]PurchasesProductRec, error) {
	var purchases []PurchasesProductRec
	err = dbmap.Raw("SELECT provider.name, purchase_detail.price, purchase.date"+
		" FROM provider, purchase_detail, purchase WHERE purchase.date>=? AND"+
		" purchase.date<= ? AND purchase_detail.product_id=? AND purchase.id"+
		"=purchase_detail.purchase_id AND provider.rut=purchase.provider_id"+
		" ORDER BY purchase.date DESC",
		in.Start, in.End, id).Scan(&purchases).Error
	return purchases, err
}
