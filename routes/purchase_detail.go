package routes

import (
	"net/http"
	"strconv"

	"coimco_backend/model"
	"github.com/gin-gonic/gin"
)

func GetPurchaseDetail(c *gin.Context) {
	purchase_id := c.Param("purchase_id")
	product_id := c.Param("product_id")
	if checkSize(purchase_id) || checkSize(product_id) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	pur_id, _ := strconv.ParseUint(purchase_id, 10, 32)
	pro_id, _ := strconv.ParseUint(product_id, 10, 32)
	purchase_detail, err := model.GetPurchaseDetail(uint(pur_id), uint(pro_id))
	if err != nil {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": GetMessageErrorSingular + " purchase_detail",
		}
		c.JSON(http.StatusNotFound, response)
	} else {
		response := gin.H{
			"status":  "success",
			"data":    purchase_detail,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	}
}

func PostPurchaseDetail(c *gin.Context) {
	var in model.PurchaseDetail
	err := c.BindJSON(&in)
	checkErr(err, BindJson)
	if err != nil || !model.CheckInPurchaseDetail(in) {
		response := gin.H{
			"status":  "error",
			"data":    nil,
			"message": ErrorParams,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	purchase_detail, flag := model.InsertPurchaseDetail(&in)

	if flag {
		response := gin.H{
			"status":  "success",
			"data":    purchase_detail,
			"message": nil,
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := gin.H{
			"status":  "error",
			"data":    purchase_detail,
			"message": PostMessageError + " a purchase_detail",
		}
		c.JSON(http.StatusBadRequest, response)
	}

}
