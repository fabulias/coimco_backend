package model

//GetSalesID return sales from that user ID
func GetSalesID(mail string, in Date) (TotalSales, error) {
	var res TotalSales
	err = dbmap.Raw("SELECT count(sale.user_id), sum(sale_detail.price*"+
		"sale_detail.quantity) FROM sale, sale_detail WHERE sale.user_id=? "+
		"AND sale.date>=? AND sale.date<=? AND sale_detail.sale_id=sale.id",
		mail, in.Start, in.End).Scan(&res).Error
	return res, err
}

//GetSales return sales in a date range
func GetSales(in Date) (TotalSales, error) {
	var res TotalSales
	err = dbmap.Raw("SELECT count(*), sum(sale_detail.price*sale_detail.quantity)"+
		" FROM sale, sale_detail WHERE sale.date>=? AND sale.date<=? "+
		"AND sale_detail.sale_id=sale.id", in.Start, in.End).Scan(&res).Error
	return res, err
}
