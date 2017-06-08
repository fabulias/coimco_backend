package model

func InsertTagCustomer(in *TagCustomer) (*TagCustomer, bool) {
	err = dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
