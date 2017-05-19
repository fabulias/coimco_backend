package model

import "strconv"

var (
	selectOneFailed = "Error selecting one row"
	selectFailed    = "Error selecting rows"
	countFailed     = "Error in select count"
)

//This function allow obtain customers' resource.
func GetCustomers(limit, offset string) ([]Customer, string) {
	var customers []Customer
	var count int64
	//Here obtain total length of table.
	count, err = dbmap.SelectInt("select count(*) from customer")
	checkErr(err, countFailed)
	//Here obtain the customers previously selected.
	_, err = dbmap.Select(&customers, "select * from customer limit $1 offset $2", limit, offset)
	checkErr(err, selectFailed)
	return customers, strconv.Itoa(int(count))
}

//This function allow obtain customer' resource for his id.
func GetCustomer(customer *Customer) *Customer {
	err := dbmap.SelectOne(&customer, "select * from customer where rut=$1", customer.Rut)
	checkErr(err, selectOneFailed)
	return customer
}

//This function allow insert customer' resource
func InsertCustomers(in *Customer) (*Customer, bool) {
	err = dbmap.Insert(in)
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
