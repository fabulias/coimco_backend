package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/fabulias/coimco_backend/model"
)

//GetInformationDashboard make route to dashboard model
func GetInformationDashboard(c *gin.Context) {
	role := c.Param("role")
	id := c.Param("id_seller")
	if role == "" {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
	} else {
		products, err1, err2 := model.GetInformationDashboard(role, id)
		if err1 != nil || err2 != nil {
			response := gin.H{
				"status":   "error",
				"data":     nil,
				"message1": DashBoardErrFirst,
				"message2": DashBoardErrSecond,
			}
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := gin.H{
				"status":  "success",
				"data":    products,
				"message": nil,
			}
			c.JSON(http.StatusOK, response)
		}
	}
}
