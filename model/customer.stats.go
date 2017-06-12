package model

func GetRankCustomerK(k string, in Date) ([]InfoCustomer, error) {
	var customers []InfoCustomer
	err = dbmap.Raw("SELECT customer.name, COUNT(sale.customer_id) AS count,"+
		" SUM(sale_detail.quantity*sale_detail.price) AS cash FROM customer, sale,"+
		" sale_detail WHERE sale.date>=? AND sale.date<=? AND customer.rut=sale."+
		"customer_id AND sale_detail.sale_id=sale.id GROUP BY customer.name ORDER "+
		"BY cash DESC LIMIT ?", in.Start, in.End, k).Scan(&customers).Error
	return customers, err
}
