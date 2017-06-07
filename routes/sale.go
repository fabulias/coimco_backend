package routes

import (
	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"net/http"
)

func GetSale(c *gin.Context) {
	customer_id := c.Param("cus_id")
	user_id := c.Param("user_id")
	sale, err := model.GetSale(customer_id, user_id)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " sale with those ID",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    sale,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

//This route insert a sale in his table
func PostSale(c *gin.Context) {
	var in model.Sale
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	//Check if client parameters are valid
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input sale
	sale, flag := model.InsertSale(&in)
	//Flag is true if the model succeeds in inserting the client
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    sale,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    sale,
			"message": PostMessageError + " a client",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}
