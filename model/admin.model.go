package model

//This function allow insert account
func InsertAccount(in *UserAcc) (*UserAcc, error) {
	err = dbmap.Create(in).Error
	return in, err
}

//This function allow obtain account' resource for his mail.
func GetAccount(mail string) (UserAcc, error) {
	var account UserAcc
	account.Mail = mail
	err = dbmap.First(&account, account.Rut).Error
	return account, err
}
