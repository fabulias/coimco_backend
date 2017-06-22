package model

//This function allow obtain customers' resource.
func GetCustomers() []Customer {
	var customers []Customer
	err = dbmap.Find(&customers).Error
	checkErr(err, selectFailed)
	return customers
}

//This function allow obtain customer' resource for his id.
func GetCustomer(rut string) (Customer, error) {
	var customer Customer
	customer.Rut = rut
	err := dbmap.Where("rut=?", rut).First(&customer).Error
	//err := dbmap.First(&customer, customer.Rut).Error
	//log.Println(err.Error())
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
