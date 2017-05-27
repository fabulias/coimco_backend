package model

type Tag struct {
	Id   int    `json:"id", db:"name:id, primarykey, autoincrement, index:idx_foreign_key_id" binding:"required"`
	Name string `json:"name", db:"name:name" binding:"required"`
}
