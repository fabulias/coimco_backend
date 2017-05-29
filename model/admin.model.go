package model

//This function allow insert account
func InsertAccount(in *User_acc) (*User_acc, bool) {
	dbmap.Create(in)
	flag := dbmap.NewRecord(in)
	return in, !flag
}

//This function allow obtain account' resource for his mail.
func GetAccount(mail string) (User_acc, error) {
	var account User_acc
	account.Mail = mail
	err = dbmap.First(&account, account.Rut).Error
	return account, err
}
