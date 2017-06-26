package routes

import (
	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"

	"net/http"
)

//GetRankSellerProductK makes route to stats model
func GetRankSellerProductK(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	err := c.BindJSON(&in)
	if err != nil || k == "" || seller == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		products, err := model.GetRankSellerProductK(k, seller, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    products,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)

		}

	}
}

//GetRankSellerProductC makes route to stats model
func GetRankSellerProductC(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	category := c.Param("category")
	err := c.BindJSON(&in)
	if err != nil || k == "" || category == "" || seller == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		products, err := model.GetRankSellerProductC(k, category, seller, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    products,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

//GetRankSellerProductB makes route to stats model
func GetRankSellerProductB(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	brand := c.Param("brand")
	err := c.BindJSON(&in)
	if err != nil || k == "" || brand == "" || seller == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		products, err := model.GetRankSellerProductB(k, brand, seller, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    products,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

//GetRankSellerCustomerK makes route to stats model
func GetRankSellerCustomerK(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	err := c.BindJSON(&in)
	if err != nil || k == "" || seller == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		customers, err := model.GetRankSellerCustomerK(k, seller, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    customers,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

//GetRankSellerCustomerP makes route to stats model
func GetRankSellerCustomerP(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	id := c.Param("id_customer")
	err := c.BindJSON(&in)
	if err != nil || k == "" || seller == "" || id == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		products, err := model.GetRankSellerCustomerP(k, id, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    products,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

//GetRankSellerCustomerL makes route to stats model
func GetRankSellerCustomerL(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	l := c.Param("l")
	err := c.BindJSON(&in)
	if err != nil || k == "" || seller == "" || l == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		customers, err := model.GetRankSellerCustomerL(k, l, seller, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    customers,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

//GetRankSellerSalesK makes route to stats model
func GetRankSellerSalesK(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	err := c.BindJSON(&in)
	if err != nil || k == "" || seller == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		customers, err := model.GetRankSellerSalesK(k, seller, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    customers,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

//GetRankSellerSalesC makes route to stats model
func GetRankSellerSalesC(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	category := c.Param("category")
	err := c.BindJSON(&in)
	if err != nil || k == "" || seller == "" || category == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		customers, err := model.GetRankSellerSalesC(k, category, seller, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    customers,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

//GetRankSellerSalesP makes route to stats model
func GetRankSellerSalesP(c *gin.Context) {
	var in model.Date
	k := c.Param("k")
	seller := c.Param("seller")
	err := c.BindJSON(&in)
	if err != nil || k == "" || seller == "" {
		resp := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, resp)
	} else {
		products, err := model.GetRankSellerSalesP(k, seller, in)
		if err != nil {
			resp := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)

		} else {
			resp := gin.H{
				"status":  "success",
				"data":    products,
				"message": nil,
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}
