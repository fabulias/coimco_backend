package routes

import (
	"coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

//This route asking for products in a range, if not exists range,
//limit and offset are 20 and 0 by default
func GetProducts(c *gin.Context) {
	//Creating limit and offset from query
	limit := c.DefaultQuery(Limit, DefaultLimit)
	offset := c.DefaultQuery(Offset, DefaultOffset)
	//Asking to model
	products, count := model.GetProducts(limit, offset)
	//Updating X-Total-Count
	c.Header(TotalCount, count)
	//If length of products is zero,
	//is because no exist products
	if checkSize(products) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorPlural + " products",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    products,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

func GetProduct(c *gin.Context) {
	mail := c.Param("mail")
	if checkSize(mail) {
		log.Println(mail)

	}
	var product *model.Product
	model.GetProduct(product)
}

//This route insert a product in his table
func PostProduct(c *gin.Context) {
	var in model.Product
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	//Check if client parameters are valid
	if !model.CheckInProduct(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": PostMessageErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input product
	product, flag := model.InsertProduct(&in)
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
