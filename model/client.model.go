package model

import "fmt"

func GetCustomers() []Client {
	var customers []Client
	_, err = dbmap.Select(&customers, "select * from customer")
	checkErr(err, "Error in Select SQL dbamp")
	fmt.Println(customers)
	return customers
}

func InsertCustomers(in Client) Client {
	err = dbmap.Insert(&in)
	checkErr(err, "Insert customer failed")
	return in
}
