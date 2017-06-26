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

//GetRankProviderVariety returns rank of provider by variety
func GetRankProviderVariety(k string, in Date) ([]ProviderRankVariety, error) {
	var providers []ProviderRankVariety
	err = dbmap.Raw("SELECT provider.name, provider.phone, provider.mail,"+
		" COUNT(purchase_detail.product_id) AS"+
		" quantity FROM provider, purchase_detail, purchase WHERE purchase.date"+
		" >= ? AND purchase.date <= ? AND provider.rut = purchase.provider_id AND"+
		" purchase_detail.purchase_id = purchase.id GROUP BY provider.name,"+
		" provider.phone, provider.mail ORDER BY quantity DESC LIMIT ?",
		in.Start, in.End, k).Scan(&providers).Error
	return providers, err
}
