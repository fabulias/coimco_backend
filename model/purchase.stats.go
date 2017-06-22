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

//GetRankPurchasesKT returns a ranking of purchase of aproducts for category
func GetRankPurchasesKT(t, k string) ([]PurchaseRankKT, error) {
	var providers []PurchaseRankKT
	err = dbmap.Raw("SELECT provider.name , date_part('day', purchase.ship_time)"+
		" AS days FROM purchase, provider WHERE date_part('day',purchase.ship_time)<?"+
		" AND provider.rut=purchase.provider_id GROUP BY provider.name,days  ORDER BY"+
		" days LIMIT ?",
		t, k).Scan(&providers).Error
	return providers, err
}
