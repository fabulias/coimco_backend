package routes

import (
	"coimco_backend/auth"
	"coimco_backend/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	var in model.Login

	if err := c.BindJSON(&in); err != nil {
		checkErr(err, BindJson)
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": BindJson,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	ret := model.LoginP(in)
	if ret {
		token, errT := auth.CreateToken(in.Mail)
		if errT != nil {
			log.Println(errT)
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": TokenError,
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response := gin.H{
			"status":  "success",
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
