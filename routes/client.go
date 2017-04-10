package routes

import (
	"coimco_backend/model"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
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
	err := c.BindJSON(&in)
	checkErr(err, "error in BindJSON")
	fmt.Println(in)

	in = model.InsertCustomers(&in)
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
