package model

//This function allow insert account
func InsertAccount(in *UserAcc) (*UserAcc, bool) {
	dbmap.Create(in)
	flag := dbmap.NewRecord(in)
	return in, !flag
}

//This function allow obtain account' resource for his mail.
func GetAccount(mail string) (UserAcc, error) {
	var account UserAcc
	account.Mail = mail
	err = dbmap.First(&account, account.Rut).Error
	return account, err
}
