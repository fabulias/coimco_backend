package model

//Customer represents the clients in the application
type Customer struct {
	Agent
}

//CustomerRankK is a struct of rank of a base
type CustomerRankK struct {
	Name  string
	Count uint
	Cash  uint
}

//CustomerRecProd is a struct for record model GetProductTotal
type CustomerRecProd struct {
	Name  string
	Total uint
}

//CustomerCash is a struct for record model GetTotalCash
type CustomerCash struct {
	Cash uint
}

type CustomerFrecuency struct {
	Name string
	Freq float64
}

type CustomerRankKL struct {
	Rut  string
	Name string
	Cant uint
}

type CustomerRankVariety struct {
	Name     string
	Quantity uint
}
