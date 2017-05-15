package model

import "log"

func GetProviders() []Provider {
	var providers []Provider
	_, err := dbmap.Select(providers, "SELECT * FROM providers")
	if err != nil {
		log.Println(providers)
		return providers
	}
	return providers
}

func InsertProviders(pin Provider) bool {
	println("Hola voy a insertar")
	errq := dbmap.Insert(&pin)
	log.Println(pin)
	log.Println(errq)
	checkErr(errq, "Insert Provider Failed")
	if errq == nil {
		println("if")
		return true
	} else {
		println("else")
		return false
	}

}
