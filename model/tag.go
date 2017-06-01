package model

import "time"

type Tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name", db:"name:name" binding:"required"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
