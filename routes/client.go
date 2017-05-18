package routes

import (
	"coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg)
	}
}

//Método que busca todos los usuarios de la bdd.
func GetCustomers(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")
	customers, count := model.GetCustomers(limit, offset)
	c.Header("X-Total-Count", count)
	if len(customers) == 0 {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": "There are no users",
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
	//mail := c.Param("mail")
	var customer *model.Customer
	model.GetCustomer(customer)
}

//Método que busca todos los usuarios de la bdd.
func PostCustomers(c *gin.Context) {
	var in model.Customer
	err := c.BindJSON(&in)
	checkErr(err, "error in BindJSON")
	if !model.CheckInCustomer(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": "I can't insert a user",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	customer, flag := model.InsertCustomers(&in)
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    customer,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "success-error",
			"data":    customer,
			"message": "Inserting client w/o arguments",
		}
		c.JSON(http.StatusOK, response)
	}

}
