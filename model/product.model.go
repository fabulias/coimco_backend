package model

//This function allow obtain products' resource.
func GetProducts() []Product {
	var products []Product
	err = dbmap.Find(&products).Error
	checkErr(err, selectFailed)
	return products
}

//This function allow obtain product' resource for his id.
func GetProduct(id uint) (Product, error) {
	var product Product
	product.ID = id
	err := dbmap.First(&product, product.ID).Error
	checkErr(err, selectOneFailed)
	return product, err
}

//This function allow insert product' resource
func InsertProduct(in *Product) (*Product, bool) {
	err = dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
