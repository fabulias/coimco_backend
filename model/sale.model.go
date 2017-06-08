package model

//GetSale insert a sale in database
func GetSale(customer_id, user_id string) (Sale, error) {
	var sale Sale
	sale.CustomerID = customer_id
	sale.UserID = user_id
	err := dbmap.Where("customer_id = $1 AND user_id = $2",
		sale.CustomerID, sale.UserID).First(&sale).Error
	checkErr(err, selectOneFailed)
	return sale, err
}

//InsertSale insert a sale in database
func InsertSale(in *Sale) (*Sale, bool) {
	err = dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
