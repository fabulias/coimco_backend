package routes

import (
	"net/http"

	"coimco_backend/model"
	"github.com/gin-gonic/gin"
)

//This route asking for providers in a range, if not exists range,
//limit and offset are 20 and 0 by default
func GetProviders(c *gin.Context) {
	//Creating limit and offset from query
	limit := c.DefaultQuery(Limit, DefaultLimit)
	offset := c.DefaultQuery(Offset, DefaultOffset)
	//Asking to model
	providers, count := model.GetProviders(limit, offset)
	//Updating X-Total-Count
	c.Header(TotalCount, count)
	//If length of providers is zero,
	//is because no exist providers
	if checkSize(providers) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorPlural + " providers",
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

//This route return a provider with a 'rut'
func GetProvider(c *gin.Context) {
	rut := c.Param("rut")
	provider, err := model.GetProvider(rut)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " provider with that rut",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    provider,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

//This route insert a provider in his table
func PostProvider(c *gin.Context) {
	var in model.Provider
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	//Check if provider parameters are valid
	if !model.CheckInProvider(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input provider
	provider, flag := model.InsertProvider(&in)
	//Flag is true if the model succeeds in inserting the provider
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    provider,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    provider,
			"message": PostMessageError + " a provider",
		}
		c.JSON(http.StatusBadRequest, response)
	}

}
