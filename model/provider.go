package model

type Provider struct {
	id    int    `db:"id, primarykey, autoincrement"`
	name  string `form:"user" json:"user" binding:"required" db:"user"`
	phone string `form:"phone" json:"phone" binding:"required" db:"phone"`
}
