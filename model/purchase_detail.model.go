package model

func GetPurchaseDetail(purchase_id, product_id uint) (PurchaseDetail, error) {
	var purchase_detail PurchaseDetail
	purchase_detail.PurchaseID = purchase_id
	purchase_detail.ProductID = product_id
	err = Dbmap.Where("purchase_id = $1 AND product_id = $2",
		purchase_detail.PurchaseID,
		purchase_detail.ProductID).First(&purchase_detail).Error
	checkErr(err, selectOneFailed)
	return purchase_detail, err
}

func InsertPurchaseDetail(in *PurchaseDetail) (*PurchaseDetail, bool) {
	err = Dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
