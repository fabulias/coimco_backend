package model

import "log"

func Products() []Product {
	var providers []Product
	_, err = dbmap.Select(providers, "SELECT * FROM providers")
	log.Println(providers)
	return 0
}
