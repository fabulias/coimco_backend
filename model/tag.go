package model

import "github.com/jinzhu/gorm"

//Struct of tag
type Tag struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
}
