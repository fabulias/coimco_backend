package routes

import (
	"net/http"

	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//This route asking for all customers
func GetCustomers(c *gin.Context) {
	//Asking to model
	customers := model.GetCustomers()
	//If length of customers is zero,
	//is because no exist customers
	if checkSize(customers) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorPlural + " clients",
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

//This route return a client with a 'rut'
func GetCustomer(c *gin.Context) {
	rut := c.Param("rut")
	customer, err := model.GetCustomer(rut)
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " client with that rut",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    customer,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

//This route insert a customer in his table
func PostCustomer(c *gin.Context) {
	var in model.Customer
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	//Check if client parameters are valid
	if !model.CheckInCustomer(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input customer
	customer, flag := model.InsertCustomer(&in)
	//Flag is true if the model succeeds in inserting the client
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    customer,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    customer,
			"message": PostMessageError + " a client",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}

//GetRankCustomerK make route to stats model
func GetRankCustomerK(c *gin.Context) {
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
		customers, err := model.GetRankCustomerK(k, in)
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

//GetProductTotal make route to record model
func GetProductTotal(c *gin.Context) {
	id := c.Param("id_customer")
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
		products, err := model.GetProductTotal(id, in)
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

//GetTotalCash make route to record model
func GetTotalCash(c *gin.Context) {
	id := c.Param("id_customer")
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
		total_cash, err := model.GetTotalCash(id, in)
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
				"data":    total_cash,
				"message": nil,
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

//GetRankFrequency make route to record model
func GetRankFrequency(c *gin.Context) {
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
		total_cash, err := model.GetRankFrequency(k, in)
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
				"data":    total_cash,
				"message": nil,
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

//GetRankCustomerKL make route to stats model
func GetRankCustomerKL(c *gin.Context) {
	k := c.Param("k")
	l := c.Param("l")
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
		customers, err := model.GetRankCustomerKL(k, l, in)
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

//GetRankCustomerVariety make route to stats model
func GetRankCustomerVariety(c *gin.Context) {
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
		customers, err := model.GetRankCustomerVariety(k, in)
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
