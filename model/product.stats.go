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

//GetRankProductCategory returns a ranking of products for category
func GetRankProductCategory(category string, in Date) ([]InfoProduct, error) {
	var products []InfoProduct
	err = dbmap.Raw("SELECT product.id, product.name , COUNT(product.id) AS "+
		"sales FROM  sale_detail, sale, product WHERE sale.date>=? AND "+
		"sale.date<=? AND sale_detail.sale_id= sale.id AND product.id="+
		"sale_detail.product_id AND product.category=? GROUP BY "+
		"product.id ORDER BY sales DESC",
		in.Start, in.End, category).Scan(&products).Error
	return products, err
}

//GetRankProductCategory returns a ranking of products for brand
func GetRankProductBrand(brand string, in Date) ([]InfoProduct, error) {
	var products []InfoProduct
	err = dbmap.Raw("SELECT product.id, product.name, COUNT(sale_detail."+
		"product_id) AS sales ,SUM(sale_detail.quantity) AS total FROM product,"+
		" sale_detail, sale WHERE sale.date>=? AND sale.date<=?"+
		" AND sale_detail.sale_id=sale.id AND product.brand=? AND "+
		"sale_detail.product_id=product.id GROUP BY product.id ORDER BY total DESC",
		in.Start, in.End, brand).Scan(&products).Error
	return products, err
}
