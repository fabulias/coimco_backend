package model

//GetRankProvidersK returns rank of providers by shiptime
func GetRankProvidersK(k string, in Date) ([]ProviderRankK, error) {
	var providers []ProviderRankK
	err = dbmap.Raw("SELECT provider.name , AVG(date_part('day', purchase."+
		"ship_time)) AS days FROM purchase, provider WHERE purchase.date>=? AND"+
		" purchase.date<= ? AND provider.rut=purchase.provider_id GROUP BY"+
		" provider.name  ORDER BY days LIMIT ?",
		in.Start, in.End, k).Scan(&providers).Error
	return providers, err
}
