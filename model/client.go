package model

type Client struct {
	id    int    `db:"id, primarykey, autoincrement"`
	name  string `form:"name" json:"name" binding:"required" db:"name"`
	phone string `form:"phone" json:"phone" binding:"required" db:"phone"`
}
