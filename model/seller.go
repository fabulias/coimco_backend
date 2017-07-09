package model

//Seller represents the seller in the application
type Seller struct {
	UserAcc
}

//This structure that serves the models
type SellerProductRank struct {
	Name string
	Cont uint
}

//This structure that serves the models
type SellerCustomerRankK struct {
	Name string
	Cash uint
}

//This structure that serves the models
type SellerCustomerRankP struct {
	Name  string
	Total uint
}

//This structure that serves the models
type SellerCustomerRankL struct {
	Name  string
	Phone string
	Mail  string
}

//This structure that serves the models
type SellerProductRec struct {
	Name  string
	Total uint
}

//This structure that serves the models
type SellerSaleRank struct {
	Name string
	Cash uint
}
