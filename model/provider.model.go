package model

import "strconv"

//This function allow obtain providers' resource.
func GetProviders(limit, offset string) ([]Provider, string) {
	var providers []Provider
	var count int64
	//Here obtain total length of table.
	count, err = dbmap.SelectInt("select count(*) from provider")
	checkErr(err, countFailed)
	//Here obtain the providers previously selected.
	_, err = dbmap.Select(&providers, "select * from provider limit $1 offset $2", limit, offset)
	checkErr(err, selectFailed)
	return providers, strconv.Itoa(int(count))
}

//This function allow obtain provider' resource for his id.
func GetProvider(rut string) (Provider, error) {
	var provider Provider
	provider.Rut = rut
	err := dbmap.SelectOne(&provider, "select * from provider where rut=$1", provider.Rut)
	checkErr(err, selectOneFailed)
	return provider, err
}

//This function allow insert provider' resource
func InsertProvider(in *Provider) (*Provider, bool) {
	err = dbmap.Insert(in)
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
