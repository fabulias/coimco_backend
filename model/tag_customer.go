package model

import "time"

//This struct
type Tag_customer struct {
	TagID      int `json:"id_tag"`
	CustomerID int `json:"id_customer"`

	Tag      Tag      `gorm:"ForeignKey:"TagID;AssociationForeignKey:ID"`
	Customer Customer `gorm:"ForeignKey:"CustomerID;AssociationForeignKey:Rut"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
