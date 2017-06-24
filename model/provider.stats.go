package model

//GetRankProviderK returns rank of providers by shiptime
func GetRankProviderK(k string, in Date) ([]ProviderRankK, error) {
	var providers []ProviderRankK
	err = dbmap.Raw("SELECT provider.name , AVG(date_part('day', purchase."+
		"ship_time)) AS days FROM purchase, provider WHERE purchase.date>=? AND"+
		" purchase.date<= ? AND provider.rut=purchase.provider_id GROUP BY"+
		" provider.name  ORDER BY days LIMIT ?",
		in.Start, in.End, k).Scan(&providers).Error
	return providers, err
}

//GetRankProviderPP returns product and price for a provider
func GetRankProviderPP(k, id string, in Date) ([]ProviderRankPP, error) {
	var products []ProviderRankPP
	err = dbmap.Raw("SELECT product.name, purchase_detail.price FROM product,"+
		" purchase, purchase_detail WHERE purchase.date>=? AND purchase.date<= ?"+
		" AND purchase.provider_id=? AND purchase_detail.purchase_id=purchase.id"+
		" AND product.id=purchase_detail.product_id GROUP BY product.name, "+
		"purchase_detail.price ORDER BY purchase_detail.price DESC LIMIT ?",
		in.Start, in.End, id, k).Scan(&products).Error
	return products, err
}
