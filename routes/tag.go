package routes

import (
	"net/http"
	"strconv"

	"coimco_backend/model"
	"github.com/gin-gonic/gin"
)

//This route insert a product in his table
func PostTag(c *gin.Context) {
	var in model.Tag
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	//Check if client parameters are valid
	if model.CheckInTag(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	product, flag := model.InsertTag(&in)
	//Flag is true if the model succeeds in inserting the client
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    product,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    product,
			"message": PostMessageError + " a client",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}

func GetTag(c *gin.Context) {
	id := c.Param("id")
	id_str, _ := strconv.ParseUint(id, 10, 64)
	product, err := model.GetTag(uint(id_str))
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " product with that ID",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    product,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}
