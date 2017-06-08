package model

import "strconv"

func GetSalesID(mail string, in Date) (TotalSalesID, string, error) {
	var count int64
	//Here obtain total length of table.
	err = dbmap.Table("sale").Count(&count).Error
	checkErr(err, countFailed)
	//Here obtain the sale previously selected.
	var res TotalSalesID

	err = dbmap.Raw("SELECT count(sale.user_id), sum(sale_detail.price*sale_detail.quantity) FROM sale, sale_detail WHERE sale.user_id=? AND sale.date>=? AND sale.date<=? AND sale_detail.sale_id=sale.id", mail, in.Start, in.End).Scan(&res).Error

	return res, strconv.Itoa(int(count)), err
}
