package model

//Represent login input
type Login struct {
	Mail string `json:"mail" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}
