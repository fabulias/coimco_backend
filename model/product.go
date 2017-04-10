package model

type Product struct {
	id        int    `db:"id, primarykey, autoincrement"`
	name      string `form:"name" json:"name" binding:"required" db:"name"`
	prototype string `form:"prototype" json:"prototype" binding:"required" db:"prototype"`
}
