package model

//GetRankCustomerK return a rank of customers in base a cash
func GetRankCustomerK(k string, in Date) ([]CustomerRankK, error) {
	var customers []CustomerRankK
	err = dbmap.Raw("SELECT customer.name, COUNT(sale.customer_id) AS count,"+
		" SUM(sale_detail.quantity*sale_detail.price) AS cash FROM customer, sale,"+
		" sale_detail WHERE sale.date>=? AND sale.date<=? AND customer.rut=sale."+
		"customer_id AND sale_detail.sale_id=sale.id GROUP BY customer.name ORDER "+
		"BY cash DESC LIMIT ?", in.Start, in.End, k).Scan(&customers).Error
	return customers, err
}

//GetRankCustomerKL return a rank of K customers of L top products
func GetRankCustomerKL(k, l string, in Date) ([]CustomerRankKL, error) {
	var customers []CustomerRankKL
	err = dbmap.Raw("SELECT customer.rut, customer.name, SUM(sale_detail."+
		"quantity) AS cant  FROM customer, sale, sale_detail, (SELECT COUNT("+
		"sale_detail.product_id) AS cantidad, sale_detail.product_id FROM "+
		"sale_detail, sale WHERE sale.date >= ? AND sale.date <= ? AND sale_"+
		"detail.sale_id=sale.id GROUP BY sale_detail.product_id ORDER BY cantidad"+
		" DESC LIMIT ?) AS products WHERE sale_detail.product_id=products."+
		"product_id AND sale.id=sale_detail.sale_id AND sale.date>=? AND"+
		" sale.date<=? AND customer.rut=sale.customer_id GROUP BY customer.rut"+
		" ORDER BY cant DESC LIMIT ?",
		in.Start, in.End, l, in.Start, in.End, k).Scan(&customers).Error
	return customers, err
}

//GetRankCustomerVariety return a rank of K customers of L top products
func GetRankCustomerVariety(k string, in Date) ([]CustomerRankVariety, error) {
	var customers []CustomerRankVariety
	err = dbmap.Raw("SELECT customer.name, COUNT(sale_detail.product_id) as"+
		" quantity FROM customer, sale_detail, sale WHERE sale.date>=? AND "+
		"sale.date<=? AND customer.rut=sale.customer_id AND sale_detail.sale_id"+
		"=sale.id GROUP BY customer.name ORDER BY quantity DESC LIMIT ?",
		in.Start, in.End, k).Scan(&customers).Error
	return customers, err
}
