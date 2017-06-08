package model

//GetPurchase insert a purchase in database
func GetPurchase(provider_id string) (Purchase, error) {
	var purchase Purchase
	purchase.ProviderID = provider_id
	err := Dbmap.Where("provider_id = $1",
		purchase.ProviderID).First(&purchase).Error
	checkErr(err, selectOneFailed)
	return purchase, err
}

//InsertPurchase insert a purchase in database
func InsertPurchase(in *Purchase) (*Purchase, bool) {
	err = Dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
