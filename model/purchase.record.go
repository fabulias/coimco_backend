package model

import "strconv"
import "log"

func GetPurcID(mail string, in Date) (TotalPurchasesID, string, error) {
	var count int64
	//Here obtain total length of table.
	log.Println(mail, in)
	err = dbmap.Table("purchase").Count(&count).Error
	checkErr(err, countFailed)
	log.Println("count -> ", count)
	//Here obtain the purchase previously selected.
	var res TotalPurchasesID

	err = dbmap.Raw("SELECT count(purchase.user_id), sum(purchase_detail.price*purchase_detail.quantity) FROM purchase, purchase_detail WHERE purchase.user_id=? AND purchase.date>=? AND purchase.date<=? AND purchase_detail.purchase_id=purchase.id", mail, in.Start, in.End).Scan(&res).Error
	log.Println("res -> ", res)
	return res, strconv.Itoa(int(count)), err
}
