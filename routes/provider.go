package routes

import (
	"coimco_backend/model"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Método que busca todos los usuarios de la bdd.
func GetProviders(c *gin.Context) {
	providers := model.GetProviders()
	if len(providers) == 0 {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": "There are no providers",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    providers,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

func PostProvider(c *gin.Context) {
	var pin model.Provider
	err := c.BindJSON(&pin)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": "Mising some field",
		}
		c.JSON(http.StatusBadRequest, response)
	} else {
		status := model.InsertProviders(pin)
		fmt.Println(pin)
		if status {
			response := gin.H{
				"status":  "success",
				"data":    nil,
				"message": "Insert Success",
			}
			c.JSON(http.StatusOK, response)
		} else {
			response := gin.H{
				"status":  "Error",
				"data":    nil,
				"message": "Provider already exist",
			}
			c.JSON(http.StatusBadRequest, response)
		}

	}

}
