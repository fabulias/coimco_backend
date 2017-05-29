package model

import "strconv"

//This function allow obtain customers' resource.
func GetCustomers(limit, offset string) ([]Customer, string) {
	var customers []Customer
	var count int64
	//Here obtain total length of table.
	err = dbmap.Table("customers").Count(count).Error
	checkErr(err, countFailed)
	//Here obtain the customers previously selected.
	err = dbmap.Offset(offset).Limit(limit).Find(&customers).Error
	checkErr(err, selectFailed)
	return customers, strconv.Itoa(int(count))
}

//This function allow obtain customer' resource for his id.
func GetCustomer(rut string) (Customer, error) {
	var customer Customer
	customer.Rut = rut
	err = dbmap.First(&customer, customer.Rut).Error
	checkErr(err, selectOneFailed)
	return customer, err
}

//This function allow insert customer' resource
func InsertCustomer(in *Customer) (*Customer, bool) {
	err = dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
