package model

import "log"

func GetProviders() []Provider {
	var providers []Provider
	_, err = dbmap.Select(providers, "SELECT * FROM providers")
	log.Println(providers)
	return providers
}

func InsertProviders(in Provider) bool {
	err = dbmap.Insert(&in)
	checkErr(err, "Insert Provider Failed")
	if err != nil {
		return true
	} else {
		return false
	}

}
