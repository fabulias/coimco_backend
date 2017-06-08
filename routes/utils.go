package routes

import (
	"github.com/fabulias/coimco_backend/model"
	"log"
	"strings"
)

//Check error function
func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg)
	}
}

func checkSize(sample interface{}) bool {
	var flag bool = false
	switch val := sample.(type) {
	case []model.Customer:
		if len(val) == 0 {
			flag = true
			return flag
		}
	case []model.Date:
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
