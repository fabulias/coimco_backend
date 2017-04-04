package model

import "log"

func Providers() int {
	var providers []Provider
	_, err = dbmap.Select(providers, "SELECT * FROM providers")
	log.Println(providers)
	return 0
}
