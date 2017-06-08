package model

func GetRankProductK(k string, in Date) ([]Product, error) {
	var products []Product
	err := dbmap.Find(&products).Error
	return products, err
}
