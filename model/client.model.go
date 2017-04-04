package model

import "log"

func Customers() []Client {
	var customers []Client
	_, err = dbmap.Select(customers, "SELECT * FROM customers")
	log.Println(customers)
	return customers
}
