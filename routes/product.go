package routes

import (
	"coimco_backend/model"
	"github.com/gin-gonic/gin"
	// _ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

//Método que busca todos los usuarios de la bdd.
func GetProducts(c *gin.Context) {
	products := model.Products()
	if len(products) == 0 {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": "There are no users",
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
