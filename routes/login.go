package routes

import (
	"net/http"

	"coimco_backend/auth"
	"coimco_backend/model"
	"github.com/gin-gonic/gin"
)

//This route generates the logic to enter in the application
func Login(c *gin.Context) {
	var in model.Login
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": BindJson,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//Check if 'in' exist in accounts with that
	//mail and pass
	acc, ret := model.LoginP(in)
	if ret {
		//Generate the token for this account
		token, errT := auth.CreateToken(in.Mail)
		if errT != nil {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": TokenError,
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		//Account information
		data := gin.H{
			"name":     acc.Name,
			"lastname": acc.Lastname,
			"role":     acc.Role,
		}
		response := gin.H{
			"status":  "success",
			"data":    data,
			"token":   token,
			"message": LoginOK,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": LoginError,
		}
		c.JSON(http.StatusBadRequest, response)
	}
}
