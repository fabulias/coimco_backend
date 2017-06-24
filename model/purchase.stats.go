package model

//GetRankPurchasesCP returns a ranking of purchase of aproducts for category
func GetRankPurchasesCP(category, k string, in Date) ([]ProductRankCategory, error) {
	var products []ProductRankCategory
	err = dbmap.Raw("SELECT products.id, products.name, SUM(purchase_detail.quantity*"+
		"purchase_detail.price) AS total FROM (SELECT * FROM product WHERE "+
		"category=? ) AS products, purchase, purchase_detail WHERE"+
		" purchase.date>=? AND purchase.date<=? AND purchase_detail.purchase_id"+
		"=purchase.id AND products.id=purchase_detail.product_id GROUP BY"+
		" products.id, products.name ORDER BY total DESC LIMIT ?",
		category, in.Start, in.End, k).Scan(&products).Error
	return products, err
}

//GetRankPurchasesK return a rank of purchases in base a total cash
func GetRankPurchasesK(k string, in Date) ([]PurchaseRankK, error) {
	var purchases []PurchaseRankK
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
