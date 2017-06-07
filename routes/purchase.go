package routes

import (
	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"net/http"
)

func GetPurchase(c *gin.Context) {
	provider_id := c.Param("prov_id")
	purchase, err := model.GetPurchase(provider_id)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " purchase with those ID",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    purchase,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

//This route insert a purchase in his table
func PostPurchase(c *gin.Context) {
	var in model.Purchase
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
	//to insert input purchase
	purchase, flag := model.InsertPurchase(&in)
	//Flag is true if the model succeeds in inserting the client
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    purchase,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    purchase,
			"message": PostMessageError + " a client",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}
