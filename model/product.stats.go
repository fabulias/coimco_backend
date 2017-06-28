package model

//GetRankProductK returns a ranking of products
func GetRankProductK(k string, in Date) ([]Product, error) {
	var products []Product
	err = dbmap.Raw("SELECT product.* FROM product, (SELECT SUM(sale_detail."+
		"quantity) AS cant, sale_detail.product_id FROM sale_detail, sale WHERE "+
		"sale.date>=? AND sale.date<=? AND "+
		"sale_detail.sale_id=sale.id GROUP BY sale_detail.product_id "+
		"ORDER BY cant DESC ) AS cant_prod WHERE "+
		"product.id=cant_prod.product_id LIMIT ?", in.Start, in.End, k).Scan(&products).Error
	return products, err
}

//GetRankProductCategoryS returns a ranking of sale products by category
func GetRankProductCategoryS(category, k string, in Date) ([]ProductRankCategory, error) {
	var products []ProductRankCategory
	err = dbmap.Raw("SELECT product.id, product.name , COUNT(product.id) AS "+
		"total FROM  sale_detail, sale, product WHERE sale.date>=? AND "+
		"sale.date<=? AND sale_detail.sale_id= sale.id AND product.id="+
		"sale_detail.product_id AND product.category=? GROUP BY "+
		"product.id ORDER BY total DESC"+
		" LIMIT ?",
		in.Start, in.End, category, k).Scan(&products).Error
	return products, err
}

//GetRankProductCategoryP returns a ranking of purchase products by category
func GetRankProductCategoryP(category, k string, in Date) ([]ProductRankCategory, error) {
	var products []ProductRankCategory
	err = dbmap.Raw("SELECT product.id, product.name , SUM(purchase_detail.quantity*purchase_detail.price) AS "+
		"total FROM  purchase_detail, purchase, product, (SELECT * FROM product WHERE category=?) AS products "+
		"WHERE purchase.date>=? AND purchase.date<=? AND purchase_detail.purchase_id= purchase.id AND products.id="+
		"purchase_detail.product_id GROUP BY product.id ORDER BY total DESC"+
		" LIMIT ?",
		category, in.Start, in.End, k).Scan(&products).Error
	return products, err
}

//GetRankProductBrand returns a ranking of products by brand
func GetRankProductBrand(brand, k string, in Date) ([]InfoProduct, error) {
	var products []InfoProduct
	err = dbmap.Raw("SELECT product.id, product.name, COUNT(sale_detail."+
		"product_id) AS sales ,SUM(sale_detail.quantity) AS total FROM product,"+
		" sale_detail, sale WHERE sale.date>=? AND sale.date<=?"+
		" AND sale_detail.sale_id=sale.id AND product.brand=? AND "+
		"sale_detail.product_id=product.id GROUP BY product.id ORDER BY total DESC"+
		" LIMIT ?",
		in.Start, in.End, brand, k).Scan(&products).Error
	return products, err
}

//GetRankProfitability returns a ranking of products by profitability
func GetRankProfitability(k string, in Date) ([]ProductRankProfitability, error) {
	var products []ProductRankProfitability
	err = dbmap.Raw("SELECT avg_product.sale-avg_product.purchase AS rent,"+
		" avg_product.name  FROM (SELECT AVG(purchase_detail.price) AS purchase,"+
		" AVG(sale_detail.price) AS sale , product.id, product.name FROM"+
		" product, sale_detail, purchase_detail, purchase, sale WHERE"+
		" purchase.date>=? AND purchase.date<=? AND purchase_detail.purchase_id"+
		" = purchase.id AND sale.date >= ? AND sale.date<= ? AND"+
		" sale_detail.sale_id=sale.id AND purchase_detail.purchase_id=purchase.id"+
		" AND sale_detail.product_id=product.id AND purchase_detail.product_id"+
		"=product.id GROUP BY product.name, product.id) AS avg_product"+
		" ORDER BY rent DESC LIMIT ?",
		in.Start, in.End, in.Start, in.End, k).Scan(&products).Error
	return products, err
}

//GetRankProductPP returns a ranking of products by provider and its price
func GetRankProductPP(id string, in Date) ([]ProductRankProviderPrice, error) {
	var products []ProductRankProviderPrice
	err = dbmap.Raw("SELECT provider.name, purchase_detail.price FROM provider,"+
		" purchase, purchase_detail WHERE purchase.date>=? AND purchase.date<=?"+
		" AND purchase_detail.product_id=? AND purchase.id="+
		"purchase_detail.purchase_id AND provider.rut=purchase.provider_id GROUP"+
		" BY provider.name, purchase_detail.price ORDER BY purchase_detail.price DESC",
		in.Start, in.End, id).Scan(&products).Error
	return products, err
}
