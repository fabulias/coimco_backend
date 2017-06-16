package model

//Customer represents the clients in the application
type Customer struct {
	Agent
}

//InfoCustomer is a struct of rank of a base
type InfoCustomer struct {
	Name  string
	Count uint
	Cash  uint
}

//CustomerRecProd is a struct for record model GetProductTotal
type CustomerRecProd struct {
	Name  string
	Total uint
}

//CustomerRecProd is a struct for record model GetProductTotal
type CustomerCash struct {
	Cash uint
}
