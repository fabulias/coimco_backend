package model

//GetRankSellerProductK returns a ranking of sold products by a seller (quantity)
func GetRankSellerProductK(k, seller string, in Date) ([]SellerProductRank, error) {
	var products []SellerProductRank
	err = dbmap.Raw("SELECT product.name, COUNT(sale_detail.product_id) AS "+
		"cont FROM sale, sale_detail, product WHERE sale.user_id=? AND "+
		"sale.date>=? AND sale.date<=? AND sale_detail.sale_id=sale.id AND"+
		" product.id=sale_detail.product_id GROUP BY product.name ORDER BY"+
		" cont DESC LIMIT ?", seller, in.Start, in.End, k).Scan(&products).Error
	return products, err

}

//GetRankSellerProductC returns a ranking of sold products by a seller and category
func GetRankSellerProductC(k, category, seller string, in Date) ([]SellerProductRank, error) {
	var products []SellerProductRank
	err = dbmap.Raw("SELECT product.name, COUNT(sale_detail.product_id) AS "+
		"cont FROM sale_detail, sale, product WHERE product.category=? AND "+
		"sale.user_id=? AND sale.date>=? AND sale.date<=? AND sale_detail.sale_id"+
		"=sale.id AND sale_detail.product_id= product.id GROUP BY product.name"+
		" ORDER BY cont DESC LIMIT ?", category,
		seller, in.Start, in.End, k).Scan(&products).Error
	return products, err

}

//GetRankSellerProductB returns a ranking of sold products by a seller and brand
func GetRankSellerProductB(k, brand, seller string, in Date) ([]SellerProductRank, error) {
	var products []SellerProductRank
	err = dbmap.Raw("SELECT product.name, COUNT(sale_detail.product_id) AS cont"+
		" FROM product, sale, sale_detail WHERE product.brand=? AND"+
		" sale.user_id=? AND sale.date>=? AND sale.date<=? AND sale_detail.sale_id"+
		"=sale.id AND sale_detail.product_id=product.id GROUP BY product.name"+
		" ORDER BY cont DESC LIMIT ?", brand, seller, in.Start, in.End, k).Scan(&products).Error
	return products, err
}

//GetRankSellerCustomerK returns a ranking of customers by a seller
func GetRankSellerCustomerK(k, seller string, in Date) ([]SellerCustomerRankK, error) {
	var customers []SellerCustomerRankK
	err = dbmap.Raw("SELECT customer.name, SUM(sale_detail.quantity*"+
		"sale_detail.price) AS cash FROM customer, sale, sale_detail WHERE"+
		" sale.user_id=? AND sale.date>=? AND sale.date<=? AND sale_detail.sale_id"+
		" = sale.id AND customer.rut=sale.customer_id GROUP BY customer.name"+
		" ORDER BY cash DESC LIMIT ?", seller, in.Start, in.End, k).Scan(&customers).Error
	return customers, err

}

//GetRankSellerCustomerP returns sold products
//(quantity) between a customer and a seller
func GetRankSellerCustomerP(k, seller, id string, in Date) ([]SellerProductRec, error) {
	var products []SellerProductRec
	err = dbmap.Raw("SELECT product.name, SUM(sale_detail.quantity) as total"+
		" FROM product, sale_detail, sale WHERE sale.user_id=? AND sale.date >=?"+
		" AND sale.date <= ? AND sale.customer_id= ? AND sale_detail.sale_id="+
		"sale.id AND product.id=sale_detail.product_id GROUP BY product.name "+
		"ORDER BY total DESC LIMIT ?", seller,
		in.Start, in.End, id, k).Scan(&products).Error
	return products, err
}

//GetRankSellerCustomerL return customers
//who do not buy the best-selling products from a seller
func GetRankSellerCustomerL(k, l, seller string, in Date) ([]SellerCustomerRankL, error) {
	var customers []SellerCustomerRankL
	err = dbmap.Raw("SELECT customer.name, customer.phone, customer.mail FROM"+
		" customer LEFT JOIN (SELECT sale.id, sale.customer_id FROM sale_detail,"+
		" sale, ( SELECT SUM(sale_detail.quantity) as cant, sale_detail."+
		"product_id FROM sale_detail, sale WHERE sale.user_id=? AND sale.date"+
		" >=? AND sale.date <= ? AND sale_detail.sale_id=sale.id GROUP BY"+
		" product_id ORDER BY cant DESC LIMIT ?) AS most_sales WHERE"+
		" sale_detail.product_id=most_sales.product_id AND sale.id="+
		"sale_detail.sale_id ) AS customer_product ON customer.rut="+
		"customer_product.customer_id GROUP BY customer.name,"+
		" customer.phone, customer.mail  LIMIT ?",
		seller, in.Start, in.End, l, k).Scan(&customers).Error
	return customers, err
}

//GetRankSellerSalesK returns a ranking of sales by a seller
func GetRankSellerSalesK(k, seller string, in Date) ([]SellerSaleRank, error) {
	var customers []SellerSaleRank
	err = dbmap.Raw("SELECT SUM(sale_detail.price*sale_detail.quantity) AS cash,"+
		" customer.name FROM customer, sale, sale_detail WHERE sale.user_id = ?"+
		" AND sale.date >=? AND sale.date<=? AND customer.rut= sale.customer_id"+
		" AND sale_detail.sale_id = sale.id GROUP BY customer.name ORDER BY"+
		" cash DESC LIMIT ?", seller, in.Start, in.End, k).Scan(&customers).Error
	return customers, err
}

//GetRankSellerSalesC returns a ranking of sales by a seller and category
func GetRankSellerSalesC(k, category, seller string, in Date) ([]SellerSaleRank, error) {
	var customers []SellerSaleRank
	err = dbmap.Raw("SELECT SUM(sale_detail.price*sale_detail.quantity) AS"+
		" cash, customer.name FROM product, customer, sale, sale_detail WHERE "+
		"sale.user_id = ? AND sale.date >=? AND sale.date<=? AND customer.rut= "+
		"sale.customer_id AND sale_detail.sale_id = sale.id AND product.category=?"+
		" AND sale_detail.product_id=product.id GROUP BY customer.name ORDER"+
		" BY cash DESC LIMIT ?", seller, in.Start, in.End, category, k).Scan(&customers).Error
	return customers, err
}

//GetRankSellerSalesP returns a ranking of sold products by a seller (money)
func GetRankSellerSalesP(k, seller string, in Date) ([]SellerSaleRank, error) {
	var products []SellerSaleRank
	err = dbmap.Raw("SELECT SUM(sale_detail.price*sale_detail.quantity) AS"+
		" cash, product.name FROM sale, sale_detail, product WHERE sale.user_id=?"+
		" AND sale.date >= ? AND sale.date <= ? AND sale_detail.sale_id= sale.id"+
		" AND product.id=sale_detail.product_id GROUP BY product.name"+
		" ORDER BY cash DESC LIMIT ?", seller, in.Start, in.End, k).Scan(&products).Error
	return products, err
}
