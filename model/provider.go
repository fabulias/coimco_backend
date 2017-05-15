package model

type Provider struct {
	Id    int64  `db:"id, primarykey, autoincrement"`
	Name  string `form:"user" json:"user" binding:"required" db:"user"`
	Phone string `form:"phone" json:"phone" binding:"required" db:"phone"`
}
