package routes

import (
	"coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

//Método que busca todos los usuarios de la bdd.
func GetProviders(c *gin.Context) {
	providers := model.Providers()
	if providers == 0 {
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
