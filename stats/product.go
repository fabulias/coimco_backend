package stats

import "github.com/fabulias/coimco_backend/model"

func GetRankProductK(k string, in model.Date) ([]model.Product, error) {
	var products []model.Product
	err := model.Dbmap.Find(&products).Error
	return products, err
}
