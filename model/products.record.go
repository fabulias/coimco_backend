package model

//GetSalesProductIDRec returns sales from product ID
func GetSalesProductIDRec(id string, in Date) ([]ProductPriceID, error) {
	var sales []ProductPriceID
	err = dbmap.Raw("SELECT SUM(sale_detail.quantity) as total, sale.date FROM "+
		"sale, sale_detail WHERE sale.date>=? AND sale.date<=? AND sale_detail."+
		"sale_id=sale.id AND sale_detail.product_id= ? GROUP BY sale.date",
		in.Start, in.End, id).Scan(&sales).Error
	return sales, err
}

//GetProductPrice returns price from product in range time
func GetProductPrice(id string, in Date) ([]ProductPrice, error) {
	var sales []ProductPrice
	err = dbmap.Raw("SELECT sale_detail.price, sale.date FROM "+
		"sale_detail, sale WHERE sale_detail.product_id=? AND sale.date>="+
		"? AND sale.date<=? AND sale.id=sale_detail.sale_id"+
		" GROUP BY sale_detail.price, sale.date",
		id, in.Start, in.End).Scan(&sales).Error
	return sales, err
}
