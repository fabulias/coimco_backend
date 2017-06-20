package model

//GetRankPurchasesK return a rank of purchases in base a total cash
func GetRankPurchasesK(k string, in Date) ([]ProviderRankK, error) {
	var purchases []ProviderRankK
	err = dbmap.Raw("SELECT provider.name AS provider_name, product.name AS "+
		"product_name, purchase_detail.quantity, purchase_detail.price,"+
		" SUM(purchase_detail.quantity*"+
		"purchase_detail.price) AS total ,purchase.id AS purchase_id FROM "+
		"provider, product, purchase, purchase_detail WHERE purchase.date >= ?"+
		" AND purchase.date <= ? AND purchase.provider_id = provider.rut AND"+
		" purchase_detail.purchase_id=purchase.id AND product.id="+
		"purchase_detail.product_id GROUP BY provider.name, product.name, "+
		"purchase_detail.quantity, purchase_detail.price, purchase.id ORDER"+
		" BY total LIMIT ?",
		in.Start, in.End, k).Scan(&purchases).Error
	return purchases, err
}
