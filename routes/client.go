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
	customers := model.GetCustomers()
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

//Método que busca todos los usuarios de la bdd.
func PostCustomers(c *gin.Context) {
	var in model.Client
	log.Println("AQUI")
	err := c.BindJSON(&in)
	checkErr(err, "error in BindJSON")
	log.Println("AQUI -> ", in)
	if !model.CheckInClient(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": "I can't insert a user",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}
	customer, _ := model.InsertCustomers(in)

	response := gin.H{
		"status":  "success",
		"data":    customer,
		"message": nil,
	}
	c.JSON(http.StatusOK, response)
}
