package model

//Seller represents the seller in the application
type Seller struct {
	UserAcc
}

type SellerProductRank struct {
	Name string
	Cont uint
}

type SellerCustomerRankK struct {
	Name string
	Cash uint
}

type SellerCustomerRankP struct {
	Name  string
	Total uint
}

type SellerCustomerRankL struct {
	Name  string
	Phone string
	Mail  string
}

type SellerProductRec struct {
	Name  string
	Total uint
}

type SellerSaleRank struct {
	Name string
	Cash uint
}
