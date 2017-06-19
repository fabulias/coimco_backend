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

//GetFrecuency returns the frecuency sales for a client
func GetFrecuency(id string, in Date) ([]CustomerFrecuency, error) {
	var customer_frecuency []CustomerFrecuency
	err = dbmap.Raw("SELECT sale.id, sale.date, sale.user_id, product.name, "+
		"sale_detail.quantity, sale_detail.price FROM sale, sale_detail, product,"+
		" customer WHERE customer.rut=? AND sale.customer_id=customer.rut"+
		" AND sale.date >=? AND sale.date <= ? AND "+
		"sale_detail.sale_id=sale.id AND product.id=sale_detail.product_id",
		id, in.Start, in.End).Scan(&customer_frecuency).Error
	return customer_frecuency, err
}
