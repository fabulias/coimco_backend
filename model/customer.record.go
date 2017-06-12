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
