package model

//GetRankSalesK returns a ranking of sales
func GetRankSalesK(k string, in Date) ([]SaleRankK, error) {
	var sales []SaleRankK
	err = dbmap.Raw("SELECT SUM(sale_detail.price*sale_detail.quantity) AS cash,"+
		" customer.name, sale.id FROM sale_detail,sale,customer WHERE sale.date>=?"+
		" AND sale.date<=? AND customer.rut=sale.customer_id AND "+
		"sale_detail.sale_id=sale.id GROUP BY sale.id,customer.name "+
		"ORDER BY cash DESC LIMIT ?",
		in.Start, in.End, k).Scan(&sales).Error
	return sales, err
}

//GetRankSalesCategory returns a ranking of sales by category
func GetRankSalesCategory(k, category string, in Date) ([]SaleRankCategory, error) {
	var sales []SaleRankCategory
	err = dbmap.Raw("SELECT SUM(sale_detail.price*sale_detail.quantity) AS cash,"+
		" customer.name, sale.id FROM sale_detail,sale,customer, product WHERE"+
		" sale.date>=? AND sale.date<=? AND customer.rut=sale.customer_id AND"+
		" sale_detail.sale_id=sale.id AND product.category=? AND"+
		" sale_detail.product_id=product.id GROUP BY sale.id,customer.name"+
		" ORDER BY cash DESC LIMIT ?",
		in.Start, in.End, category, k).Scan(&sales).Error
	return sales, err
}

//GetRankSalesProduct returns a ranking of sold products
func GetRankSalesProduct(k string, in Date) ([]SaleRankProduct, error) {
	var products []SaleRankProduct
	err = dbmap.Raw("SELECT SUM(sale_detail.quantity*sale_detail.price) AS cash,"+
		" product.id, product.name FROM product, sale_detail, sale WHERE"+
		" sale.date>=? AND sale.date<=? AND sale_detail.sale_id=sale.id AND"+
		" product.id=sale_detail.product_id GROUP BY product.name,"+
		" product.id ORDER BY cash DESC LIMIT ?",
		in.Start, in.End, k).Scan(&products).Error
	return products, err
}
