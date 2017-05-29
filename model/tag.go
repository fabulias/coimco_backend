package model

import "time"

type Tag struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name", db:"name:name" binding:"required"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
