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

//GetRankProductCategoryS returns a ranking of sale products for category
func GetRankProductCategoryS(category string, in Date) ([]ProductRankCategory, error) {
	var products []ProductRankCategory
	err = dbmap.Raw("SELECT product.id, product.name , COUNT(product.id) AS "+
		"total FROM  sale_detail, sale, product WHERE sale.date>=? AND "+
		"sale.date<=? AND sale_detail.sale_id= sale.id AND product.id="+
		"sale_detail.product_id AND product.category=? GROUP BY "+
		"product.id ORDER BY sales DESC",
		in.Start, in.End, category).Scan(&products).Error
	return products, err
}

//GetRankProductCategoryP returns a ranking of purchase products for category
func GetRankProductCategoryP(category string, in Date) ([]ProductRankCategory, error) {
	var products []ProductRankCategory
	err = dbmap.Raw("SELECT product.id, product.name , SUM(purchase_detail.quantity*purchse_detail.price) AS "+
		"total FROM  purchase_detail, purchase, product, (SELECT * FROM product WHERE category=?) AS products "+
		"WHERE purchase.date>=? AND purchase.date<=? AND purchase_detail.purchase_id= purchase.id AND products.id="+
		"purchase_detail.product_id GROUP BY product.id ORDER BY sales DESC",
		category, in.Start, in.End).Scan(&products).Error
	return products, err
}

//GetRankProductBrand returns a ranking of products for brand
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
