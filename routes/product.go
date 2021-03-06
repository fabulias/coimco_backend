package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/fabulias/coimco_backend/model"
)

//This route asking for all products
func GetProducts(c *gin.Context) {
	//Asking to model
	products := model.GetProducts()
	//If length of products is zero,
	//is because no exist products
	if checkSize(products) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorPlural + " products",
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

//GetProduct return information from this product
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	id_str, _ := strconv.ParseUint(id, 10, 64)
	product, err := model.GetProduct(uint(id_str))
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " product with that ID",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    product,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

//This route insert a product in his table
func PostProduct(c *gin.Context) {
	var in model.Product
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	//Check if client parameters are valid
	if !model.CheckInProduct(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//As the params are correct, we proceeded
	//to insert input product
	product, flag := model.InsertProduct(&in)
	//Flag is true if the model succeeds in inserting the client
	if flag {
		response := gin.H{
			"status":  "success",
			"data":    product,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    product,
			"message": PostMessageError + " a client",
		}
		c.JSON(http.StatusBadRequest, response)
	}
}

//GetRankProductK make route to stats model
func GetRankProductK(c *gin.Context) {
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
		products, err := model.GetRankProductK(k, in)
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

//GetSalesProductIDRec make route to model
func GetSalesProductIDRec(c *gin.Context) {
	id := c.Param("id")
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
		sales, err := model.GetSalesProductIDRec(id, in)
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
				"data":    sales,
				"message": nil,
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

//GetRankProductCategoryS make route to stats model
func GetRankProductCategoryS(c *gin.Context) {
	category := c.Param("category")
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
		products, err := model.GetRankProductCategoryS(category, k, in)
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

//GetRankProductCategoryP make route to stats model
func GetRankProductCategoryP(c *gin.Context) {
	category := c.Param("category")
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
		products, err := model.GetRankProductCategoryP(category, k, in)
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

//GetRankProductBrand make route to stats model
func GetRankProductBrand(c *gin.Context) {
	brand := c.Param("brand")
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
		products, err := model.GetRankProductBrand(brand, k, in)
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

//GetRankProfitability make route to stats model
func GetRankProfitability(c *gin.Context) {
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
		products, err := model.GetRankProfitability(k, in)
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

//GetRankProductPP make route to stats model
func GetRankProductPP(c *gin.Context) {
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
		products, err := model.GetRankProductPP(id, in)
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
