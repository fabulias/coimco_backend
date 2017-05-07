package model

import "log"

func GetCustomers() []Client {
	var customers []Client
	log.Println("ASD")
	_, err = dbmap.Select(&customers, "select * from customer")
	checkErr(err, "Error in Select SQL dbamp")
	return customers
}

func InsertCustomers(in *Client) (*Client, bool) {
	err = dbmap.Insert(in)
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
