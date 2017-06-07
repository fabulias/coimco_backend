package routes

import (
	"log"
	"net/http"

	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"
)

func PostTagCustomer(c *gin.Context) {
	var in model.TagCustomer
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	log.Println("in-> ", in)
	//Check if tag parameters are valid
	if model.CheckInTagCustomer(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	tag, flag := model.InsertTagCustomer(&in)
	//Flag is true if the model succeeds in inserting the tag
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    tag,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    tag,
			"message": PostMessageError + " a tag",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}
