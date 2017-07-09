package routes

import (
	"net/http"
	"strconv"

	"github.com/fabulias/coimco_backend/model"
	"github.com/gin-gonic/gin"
)

//GetSaleDetail makes route to model
func GetSaleDetail(c *gin.Context) {
	sale_id := c.Param("sale_id")
	product_id := c.Param("product_id")
	if checkSize(sale_id) || checkSize(product_id) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	s_id, _ := strconv.ParseUint(sale_id, 10, 32)
	pro_id, _ := strconv.ParseUint(product_id, 10, 32)
	sale_detail, err := model.GetSaleDetail(uint(s_id), uint(pro_id))
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " sale_detail",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    sale_detail,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

//PostSaleDetail makes route to model
func PostSaleDetail(c *gin.Context) {
	var in model.SaleDetail
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	if err != nil || !model.CheckInSaleDetail(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	sale_detail, flag := model.InsertSaleDetail(&in)

	if flag {
		response := gin.H{
			"status":  "success",
			"data":    sale_detail,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    sale_detail,
			"message": PostMessageError + " a sale_detail",
		}
		c.JSON(http.StatusBadRequest, response)
	}

}
