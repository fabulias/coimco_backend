package model

//This function allow insert account
func InsertAccount(in *User_acc) (*User_acc, bool) {
	err = dbmap.Insert(in)
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}

//This function allow obtain account' resource for his mail.
func GetAccount(mail string) (User_acc, error) {
	var account User_acc
	account.Mail = mail
	err := dbmap.SelectOne(&account, "select * from user_acc where mail=$1", account.Mail)
	checkErr(err, selectOneFailed)
	return account, err
}
