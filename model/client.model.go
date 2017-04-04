package model

import "log"

func Customers() int {
	var customers []Client
	_, err = dbmap.Select(customers, "SELECT * FROM customers")
	log.Println(customers)
	return 0
}
