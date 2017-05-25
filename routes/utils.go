package routes

import (
	"coimco_backend/model"
	"strings"
)

func checkSize(sample interface{}) bool {
	var flag bool = false
	switch val := sample.(type) {
	case []model.Customer:
		if len(val) == 0 {
			flag = true
			return flag
		}
	case []model.Product:
		if len(val) == 0 {
			flag = true
			return flag
		}
	case string:
		if strings.Compare(val, "") == 0 {
			flag = true
			return flag
		}
	}
	return flag

}
