package routes

import (
	"net/http"

	"coimco_backend/hash"
	"coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//This route insert an account in user_acc table
func PostAccount(c *gin.Context) {
	var in model.UserAcc
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	if model.CheckInAccount(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input account
	in.Active = true
	hash_pass, err := hash.HashPassword(in.Pass)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorHashPassword,
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	in.Pass = hash_pass
	account, err := model.InsertAccount(&in)
	//Flag is true if the model succeeds in inserting account
	if err != nil {
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
			"message": PostMessageError + " account",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}

//This route return a account with a 'mail'
func GetAccount(c *gin.Context) {
	mail := c.Param("mail")
	account, err := model.GetAccount(mail)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " account with that mail",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    account,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}
