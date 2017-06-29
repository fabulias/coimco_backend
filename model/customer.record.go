package model

//GetProductTotal returns the total sales of a product for this client
func GetProductTotal(id string, in Date) ([]CustomerRecProd, error) {
	var products []CustomerRecProd
	err = dbmap.Raw("SELECT product.name, SUM(sale_detail.quantity) AS total"+
		" FROM product, sale_detail, sale WHERE sale.customer_id=? AND sale.date>=? "+
		"AND sale.date<=? AND sale_detail.sale_id=sale.id AND product.id="+
		"sale_detail.product_id GROUP BY product.name ORDER BY total DESC",
		id, in.Start, in.End).Scan(&products).Error
	return products, err
}

//GetTotalCash returns the total sales of a product for this client
func GetTotalCash(id string, in Date) (CustomerCash, error) {
	var total_cash CustomerCash
	err = dbmap.Raw("SELECT SUM(sale_detail.quantity*sale_detail.price) AS cash"+
		" FROM sale, sale_detail, customer WHERE customer.rut=? AND "+
		"sale.customer_id=customer.rut AND sale.date >= ? AND sale.date"+
		"<= ? AND sale_detail.sale_id=sale.id",
		id, in.Start, in.End).Scan(&total_cash).Error
	return total_cash, err
}

//GetRankFrequency returns the frecuency sales for all clients
func GetRankFrequency(k string, in Date) ([]CustomerFrecuency, error) {
	duration := in.End.Sub(in.Start)
	var customer_frecuency []CustomerFrecuency
	err = dbmap.Raw("SELECT COUNT(sale.customer_id)::float/(?::float) as freq,"+
		" customer.name as name FROM sale, customer"+
		" GROUP BY customer_id ORDER BY freq DESC LIMIT ?",
		duration.Hours()/24/30, k).Scan(&customer_frecuency).Error
	return customer_frecuency, err
}
