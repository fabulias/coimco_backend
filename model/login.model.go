package model

import "coimco_backend/hash"

//This function allow sign in an account
func LoginP(in Login) (User_acc, bool) {
	user_acc, err := GetAccount(in.Mail)
	//If the account exists
	if err == nil {
		//Check password and hash
		flag := hash.CheckPasswordHash(in.Pass, user_acc.Pass)
		//Password and hash are equals
		if flag {
			return user_acc, true
		}
	}
	return user_acc, false
}
