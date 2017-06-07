package model

import "strconv"

//This function allow obtain provider' resource.
func GetProviders(limit, offset string) ([]Provider, string) {
	var provider []Provider
	var count int64
	//Here obtain total length of table.
	err = dbmap.Table("provider").Count(count).Error
	checkErr(err, countFailed)
	//Here obtain the provider previously selected.
	err = dbmap.Offset(offset).Limit(limit).Find(&provider).Error
	checkErr(err, selectFailed)
	return provider, strconv.Itoa(int(count))
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
