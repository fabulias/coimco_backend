package routes

import (
	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"log"
	"net/http"
	"strings"
)

func GetSales(c *gin.Context) {
	//Creating limit and offset from query
	mail := c.Param("mail")
	var in model.Date
	err := c.BindJSON(&in)
	log.Println("mail -> ", mail)
	log.Println("in -> ", in)
	if strings.Compare(mail, "") == 0 || err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": BindJson,
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		//Asking to model
		sales, count := model.GetSales(mail, in)
		//Updating X-Total-Count
		c.Header(TotalCount, count)
		//If length of sales is zero,
		//is because no exist sales
		if checkSize(sales) {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": GetMessageErrorPlural + " clients",
			}
			c.JSON(http.StatusNotFound, response)
		} else {
			response := gin.H{
				"status":  "success",
				"data":    sales,
				"message": nil,
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

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
			"message": PostMessageError + " a sale",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}
