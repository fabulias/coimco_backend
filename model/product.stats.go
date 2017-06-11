package model

//GetRankProductK returns a ranking of products
func GetRankProductK(k string, in Date) ([]Product, error) {
	var products []Product
	err = dbmap.Raw("SELECT product.* FROM product, (SELECT count(sale_detail."+
		"product_id) AS cant, sale_detail.product_id FROM sale_detail, sale WHERE "+
		"sale.date>=? AND sale.date<=? AND "+
		"sale_detail.sale_id=sale.id GROUP BY sale_detail.product_id "+
		"ORDER BY cant DESC ) AS cant_prod WHERE "+
		"product.id=cant_prod.product_id LIMIT ?", in.Start, in.End, k).Scan(&products).Error
	return products, err
}
