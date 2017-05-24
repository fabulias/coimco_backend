package model

func Providers() int {
	var providers []Provider
	_, err = dbmap.Select(providers, "SELECT * FROM providers")
	return 0
}
