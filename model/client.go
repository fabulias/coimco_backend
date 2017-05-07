package model

type Client struct {
	Id    int64  `db:"id, primarykey, autoincrement"`
	Name  string `form:"name" json:"name" binding:"required" db:"name"`
	Phone string `form:"phone" json:"phone" binding:"required" db:"phone"`
}
