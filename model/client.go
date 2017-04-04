package model

type Client struct {
	Id    int    `db:"id, primarykey, autoincrement"`
	Name  string `form:"user" json:"user" binding:"required" db:"user"`
	Phone string `form:"phone" json:"phone" binding:"required" db:"phone"`
}
