package model

import "strconv"
import "log"

//This function allow obtain sales' resource.
func GetSales(mail string, in Date) (Sale, string) {
	var sale Sale
	var count int64
	//Here obtain total length of table.
	err = Dbmap.Table("sale").Count(count).Error
	checkErr(err, countFailed)
	//Here obtain the sale previously selected.
	algo := Dbmap.Where("SELECT count(sale.user_id), sum(sale_detail.price) FROM sale, sale_detail WHERE sale.user_id=? AND sale.date>=? AND sale.date<=? AND sale_detail.sale_id=sale.id", mail, in.Start, in.End).Find(&sale).Value
	log.Println(algo)
	//checkErr(err, selectFailed)
	return sale, strconv.Itoa(int(count))
}

//GetSale insert a sale in database
func GetSale(customer_id, user_id string) (Sale, error) {
	var sale Sale
	sale.CustomerID = customer_id
	sale.UserID = user_id
	err := Dbmap.Where("customer_id = $1 AND user_id = $2",
		sale.CustomerID, sale.UserID).First(&sale).Error
	checkErr(err, selectOneFailed)
	return sale, err
}

//InsertSale insert a sale in database
func InsertSale(in *Sale) (*Sale, bool) {
	err = Dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
