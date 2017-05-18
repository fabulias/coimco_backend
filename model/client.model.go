package model

//import "crypto/md5"
import (
	"strconv"
)

func GetCustomers(limit, offset string) ([]Customer, string) {
	var customers []Customer
	var count int64
	count, err = dbmap.SelectInt("select count(*) from customer")
	checkErr(err, "count data return err")
	_, err = dbmap.Select(&customers, "select * from customer limit $1 offset $2", limit, offset)
	checkErr(err, "Error in Select SQL dbamp")
	return customers, strconv.Itoa(int(count))
}

func GetCustomer(customer *Customer) *Customer {
	err := dbmap.SelectOne(&customer, "select * from customer where rut=$1", customer.Rut)
	checkErr(err, "SelectOne failed")
	return customer
}

func InsertCustomers(in *Customer) (*Customer, bool) {
	err = dbmap.Insert(in)
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
