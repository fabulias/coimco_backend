package model

import "log"

func Products() []Product {
	var products []Product
	_, err = dbmap.Select(&products, "SELECT * FROM product")
	log.Println(products)
	return products
}
