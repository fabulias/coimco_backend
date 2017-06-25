package model

import "log"

//GetInformationDashboard returns sales from product ID
func GetInformationDashboard(role, id string) (DashBoardInformation, error, error) {
	var information DashBoardInformation
	var err1 error = nil
	var err2 error = nil
	//If you need information to dashboard seller
	if role == "0" {
		log.Println("Seller query")
		err1 = dbmap.Raw("SELECT cash_sales.sale_total, SUM(sale_detail.quantity"+
			"*sale_detail.price) as last_sale FROM sale, sale_detail, ( SELECT "+
			"SUM(sale_detail.quantity*sale_detail.price) as sale_total, "+
			"MAX(last_sales.id) as id FROM (SELECT sale.id, sale.date FROM sale"+
			" WHERE sale.user_id=? ORDER BY sale.date DESC LIMIT 7) AS "+
			"last_sales, sale_detail WHERE sale_detail.sale_id= last_sales.id)"+
			" AS cash_sales WHERE sale_detail.sale_id=cash_sales.id "+
			"GROUP BY sale_total", id).Scan(&information).Error
	} else {
		log.Println("Admin and manager querys")
		//This case is to admin and manager
		err1 = dbmap.Raw("SELECT cash_purchase.purchase_total, " +
			"SUM(purchase_detail.quantity*purchase_detail.price) as last_purchase" +
			" FROM ( SELECT SUM(purchase_detail.quantity*purchase_detail.price) AS" +
			" purchase_total, MAX(last_purchase.id) as id FROM (SELECT purchase.id," +
			" purchase.date FROM purchase ORDER BY purchase.date ) AS last_purchase," +
			" purchase_detail WHERE purchase_detail.purchase_id=last_purchase.id) AS" +
			" cash_purchase, purchase_detail WHERE purchase_detail.purchase_id=" +
			"cash_purchase.id GROUP BY purchase_total").Scan(&information).Error
		err2 = dbmap.Raw("SELECT cash_sale.sale_total, SUM(sale_detail.quantity" +
			"*sale_detail.price) as last_sale FROM (SELECT SUM(sale_detail.quantity" +
			"*sale_detail.price) as sale_total, MAX(last_sales.id) AS id FROM " +
			"(SELECT sale.id, sale.date FROM sale ORDER BY sale.date DESC LIMIT 7)" +
			" AS last_sales, sale_detail WHERE sale_detail.sale_id=last_sales.id )" +
			" AS cash_sale, sale_detail WHERE sale_detail.sale_id=cash_sale.id" +
			" GROUP BY cash_sale.sale_total").Scan(&information).Error
	}
	return information, err1, err2
}
