package model

func InsertAccount(in *User_acc) (*User_acc, bool) {
	err = dbmap.Insert(in)
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
