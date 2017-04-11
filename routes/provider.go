package routes

import (
	"coimco_backend/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

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

func PostProviders(c *gin.Context) {
	var in model.Provider
	err := c.BindJSON(&in)
	checkErr(err, "Error in BindJSON")

	in = model.InsertProviders(in)
	if len(providers) == 0 {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": "There are no user",
		}
		c.JSON(response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    providers,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}

}
