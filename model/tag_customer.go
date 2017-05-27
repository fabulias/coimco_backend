package model

type Tag_client struct {
	Id_tag int `json:"id_tag", db:"name:id_tag, relation:Id"`
}
