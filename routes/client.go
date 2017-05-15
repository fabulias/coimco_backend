package routes

import (
	"coimco_backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//Método que busca todos los usuarios de la bdd.
func GetCustomers(c *gin.Context) {
	customers := model.Customers()
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
