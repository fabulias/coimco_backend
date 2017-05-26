package model

import "strings"
import "log"

//Return true in case of that all params are okay
func CheckInCustomer(in Customer) bool {
	if strings.Compare(in.Name, "") != 0 && strings.Compare(in.Rut, "") != 0 && strings.Compare(in.Mail, "") != 0 {
		return true
	} else {
		return false
	}
}

//Return true in case of that all params are okay
func CheckInAccount(in User_acc) bool {
	var flag bool = false
	if strings.Compare(in.Name, "") == 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Lastname, "") == 0 {
		flag = true
		return flag
	} else if in.Role < 0 && in.Role > 2 {
		flag = true
		return flag
	} else if strings.Compare(in.Mail, "") == 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Rut, "") == 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Pass, "") == 0 {
		flag = true
		return flag
	}
	return flag
}

//Return true in case of that all params are okay
func CheckInProduct(in Product) bool {
	var flag bool = false
	if strings.Compare(in.Name, "") != 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Details, "") != 0 {
		flag = true
		return flag
	} else if in.Stock < 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Brand, "") != 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Category, "") != 0 {
		flag = true
		return flag
	}
	return flag
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg)
	}
}
