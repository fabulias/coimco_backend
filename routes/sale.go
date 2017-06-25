package routes

import (
	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"net/http"
	"strings"
)

//GetSalesID bind JSON, param URI inputs and call model stats
func GetSalesID(c *gin.Context) {
	mail := c.Param("mail")
	var in model.Date
	err := c.BindJSON(&in)
	if strings.Compare(mail, "") == 0 || err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": BindJson,
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		res, err := model.GetSalesID(mail, in)
		if err != nil {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": GetMessageErrorPlural + " sales",
			}
			c.JSON(http.StatusNotFound, response)
		} else {
			response := gin.H{
				"status":  "success",
				"data":    res,
				"message": nil,
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

//GetSales bind JSON input and call model stats
func GetSales(c *gin.Context) {
	var in model.Date
	err := c.BindJSON(&in)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": BindJson,
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		//Asking to model
		res, err := model.GetSales(in)
		//If length of sales is zero,
		//is because no exist sales
		if err != nil {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": GetMessageErrorPlural + " sales",
			}
			c.JSON(http.StatusNotFound, response)
		} else {
			response := gin.H{
				"status":  "success",
				"data":    res,
				"message": nil,
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

func GetSale(c *gin.Context) {
	customer_id := c.Param("cus_id")
	user_id := c.Param("user_id")
	sale, err := model.GetSale(customer_id, user_id)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " sale with those ID",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    sale,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

//This route insert a sale in his table
func PostSale(c *gin.Context) {
	var in model.Sale
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	//Check if client parameters are valid
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input sale
	sale, flag := model.InsertSale(&in)
	//Flag is true if the model succeeds in inserting the client
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    sale,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    sale,
			"message": PostMessageError + " a sale",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}

//GetRankSalesK makes route to stats model
func GetRankSalesK(c *gin.Context) {
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
		products, err := model.GetRankSalesK(k, in)
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

//GetRankSalesCategory makes route to stats model
func GetRankSalesCategory(c *gin.Context) {
	k := c.Param("k")
	category := c.Param("category")
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
		products, err := model.GetRankSalesCategory(k, category, in)
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

//GetRankSalesProduct makes route to stats model
func GetRankSalesProduct(c *gin.Context) {
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
		products, err := model.GetRankSalesProduct(k, in)
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

//GetRankSalesArea makes route to stats model
func GetRankSalesArea(c *gin.Context) {
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
		products, err := model.GetRankSalesArea(k, in)
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

//GetSalesProduct makes route to record model
func GetSalesProduct(c *gin.Context) {
	id := c.Param("id_product")
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
		products, err := model.GetSalesProduct(id, in)
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
