package routes

import (
	"log"
	"net/http"

	"coimco_backend/model"
	"github.com/gin-gonic/gin"
)

func PostAccount(c *gin.Context) {
	var in model.User_acc
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	log.Println(in)
	if !model.CheckInAccount(in) {
		log.Println("HERE")
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": PostMessageErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input account
	in.Active = true
	log.Println(in)
	account, flag := model.InsertAccount(&in)
	//Flag is true if the model succeeds in inserting the client
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    account,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    account,
			"message": PostMessageError + " a account",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}
