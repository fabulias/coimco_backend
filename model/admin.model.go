package model

//This function allow insert account
func InsertAccount(in *UserAcc) (*UserAcc, error) {
	err = Dbmap.Create(in).Error
	return in, err
}

//This function allow obtain account' resource for his mail.
func GetAccount(mail string) (UserAcc, error) {
	var account UserAcc
	account.Mail = mail
	err := Dbmap.Where("mail=?", mail).First(&account).Error
	return account, err
}
