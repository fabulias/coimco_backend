package model

func GetPurchaseDetail(purchase_id, product_id uint) (PurchaseDetail, error) {
	var purchase_detail PurchaseDetail
	purchase_detail.PurchaseID = purchase_id
	purchase_detail.ProductID = product_id
	err = dbmap.First(&purchase_detail, purchase_detail.PurchaseID,
		purchase_detail.ProductID).Error
	checkErr(err, selectOneFailed)
	return purchase_detail, err
}

func InsertPurchaseDetail(in *PurchaseDetail) (*PurchaseDetail, bool) {
	err = dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
