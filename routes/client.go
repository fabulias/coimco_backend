package routes

import (
	"coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

//Check error function
func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg)
	}
}

var (
	Limit         = "limit"
	Offset        = "offset"
	DefaultLimit  = "20"
	DefaultOffset = "0"
	TotalCount    = "X-Total-Count"
)

//This route asking for customers in a range, if not exists range,
//limit and offset are 20 and 0 by default
func GetCustomers(c *gin.Context) {
	//Creating limit and offset from query
	limit := c.DefaultQuery(Limit, DefaultLimit)
	offset := c.DefaultQuery(Offset, DefaultOffset)
	//Asking to model
	customers, count := model.GetCustomers(limit, offset)
	//Updating X-Total-Count
	c.Header(TotalCount, count)
	//If length of customers is zero,
	//is because no exist customers
	if checkSize(customers) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GET_MESSAGE_ERROR_PLURAL + " clients",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    customers,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

func GetCustomer(c *gin.Context) {
	mail := c.Param("mail")
	if checkSize(mail) {
		log.Println(mail)

	}
	var customer *model.Customer
	model.GetCustomer(customer)
}

//This route insert a customer in his table
func PostCustomers(c *gin.Context) {
	var in model.Customer
	err := c.BindJSON(&in)
	checkErr(err, BIND_JSON)
	//Check if client parameters are valid
	if !model.CheckInCustomer(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": POST_MESSAGE_ERROR_PARAMS,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input customer
	customer, flag := model.InsertCustomers(&in)
	//Flag is true if the model succeeds in inserting the client
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    customer,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    customer,
			"message": POST_MESSAGE_ERROR + " a client",
		}
		c.JSON(http.StatusBadRequest, response)
	}

}
