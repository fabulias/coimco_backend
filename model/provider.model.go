package model

//This function allow obtain provider' resource.
func GetProviders() []Provider {
	var provider []Provider
	err = dbmap.Find(&provider).Error
	checkErr(err, selectFailed)
	return provider
}

//This function allow obtain provider' resource for his id.
func GetProvider(rut string) (Provider, error) {
	var provider Provider
	provider.Rut = rut
	err := dbmap.Where("rut=?", rut).First(&provider).Error
	checkErr(err, selectOneFailed)
	return provider, err
}

//This function allow insert provider' resource
func InsertProvider(in *Provider) (*Provider, bool) {
	err = dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
