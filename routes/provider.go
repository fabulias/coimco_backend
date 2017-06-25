package routes

import (
	"net/http"

	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"
)

//This route asking for all providers
func GetProviders(c *gin.Context) {
	//Asking to model
	providers := model.GetProviders()
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

//GetRankPurchasesK make route to stats model
func GetRankPurchasesK(c *gin.Context) {
	k := c.Param("k")
	var in model.Date
	err := c.BindJSON(&in)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
	} else {
		customers, err := model.GetRankPurchasesK(k, in)
		if err != nil {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := gin.H{
				"status":  "success",
				"data":    customers,
				"message": nil,
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

//GetRankProviderK make route to stats model
func GetRankProviderK(c *gin.Context) {
	k := c.Param("k")
	var in model.Date
	err := c.BindJSON(&in)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
	} else {
		products, err := model.GetRankProviderK(k, in)
		if err != nil {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
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

//GetRankProviderPP make route to stats model
func GetRankProviderPP(c *gin.Context) {
	k := c.Param("k")
	id := c.Param("id_provider")
	var in model.Date
	err := c.BindJSON(&in)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
	} else {
		products, err := model.GetRankProviderPP(k, id, in)
		if err != nil {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
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

//GetRankProviderVariety make route to stats model
func GetRankProviderVariety(c *gin.Context) {
	k := c.Param("k")
	var in model.Date
	err := c.BindJSON(&in)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
	} else {
		products, err := model.GetRankProviderVariety(k, in)
		if err != nil {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
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
