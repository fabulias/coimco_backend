package model

import "strconv"

//This function allow obtain products' resource.
func GetProducts(limit, offset string) ([]Product, string) {
	var products []Product
	var count int64
	//Here obtain total length of table.
	Dbmap.Table("products").Count(count)
	//Here obtain the products previously selected.
	Dbmap.Limit(limit).Offset(offset).Find(&products)
	return products, strconv.Itoa(int(count))
}

//This function allow obtain product' resource for his id.
func GetProduct(id uint) (Product, error) {
	var product Product
	product.ID = id
	err := Dbmap.First(&product, product.ID).Error
	checkErr(err, selectOneFailed)
	return product, err
}

//This function allow insert product' resource
func InsertProduct(in *Product) (*Product, bool) {
	err = Dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
