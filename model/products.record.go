package model

//GetSalesProductIDRec returns sales from product ID
func GetSalesProductIDRec(id string, in Date) ([]Sale, error) {
	var sales []Sale
	err = dbmap.Raw("SELECT sale.customer_id, sale.user_id, sale_detail.quantity"+
		", sale.date FROM sale, sale_detail WHERE sale_detail.product_id=? AND"+
		" sale.id=sale_detail.sale_id AND sale.date>=?' AND sale.date<=?",
		in.Start, in.End, id).Scan(&sales).Error
	return sales, err
}
