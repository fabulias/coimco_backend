package model

func GetSaleDetail(purchase_id, product_id uint) (SaleDetail, error) {
	var sale_detail SaleDetail
	sale_detail.SaleID = purchase_id
	sale_detail.ProductID = product_id
	err = Dbmap.First(&sale_detail, sale_detail.SaleID,
		sale_detail.ProductID).Error
	checkErr(err, selectOneFailed)
	return sale_detail, err
}

func InsertSaleDetail(in *SaleDetail) (*SaleDetail, bool) {
	err = Dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
