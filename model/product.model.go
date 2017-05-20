package model

import "strconv"

//This function allow obtain products' resource.
func GetProducts(limit, offset string) ([]Product, string) {
	var products []Product
	var count int64
	//Here obtain total length of table.
	count, err = dbmap.SelectInt("select count(*) from product")
	checkErr(err, countFailed)
	//Here obtain the products previously selected.
	_, err = dbmap.Select(&products, "select * from product limit $1 offset $2", limit, offset)
	checkErr(err, selectFailed)
	return products, strconv.Itoa(int(count))
}

//This function allow obtain product' resource for his id.
func GetProduct(product *Product) *Product {
	err := dbmap.SelectOne(&product, "select * from product where id=$1", product.Id)
	checkErr(err, selectOneFailed)
	return product
}

//This function allow insert product' resource
func InsertProduct(in *Product) (*Product, bool) {
	err = dbmap.Insert(in)
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
