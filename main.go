package main

import (
	"log"
	"os"

	"github.com/fabulias/coimco_backend/auth"
	"github.com/fabulias/coimco_backend/routes"
	"github.com/gin-gonic/gin"
)

var StatusOK = 200

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(StatusOK)
		} else {
			c.Next()
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(Cors())
	r.POST("/login", routes.Login)
	// Simple group: v1
	v1 := r.Group("api")
	v1.Use(auth.ValidateToken())
	{
		//Methods plurals GET
		v1.GET("/customers", routes.GetCustomers)
		v1.GET("/products", routes.GetProducts)
		v1.GET("/providers", routes.GetProviders)

		//Methods singular GET
		v1.GET("/customers/:rut", routes.GetCustomer)
		v1.GET("/providers/:rut", routes.GetProvider)
		v1.GET("/products/:id", routes.GetProduct)
		v1.GET("/accounts/:mail", routes.GetAccount)
		v1.GET("/tags/:id", routes.GetTag)
		v1.GET("/sale_detail/:sale_id/:product_id", routes.GetSaleDetail)
		v1.GET("/purchase_detail/:purchase_id/:product_id", routes.GetPurchaseDetail)
		v1.GET("/sales/:cus_id/:user_id", routes.GetSale)
		v1.GET("/purchases/:prov_id", routes.PostPurchase)

		//Methods POST
		v1.POST("/customers", routes.PostCustomer)
		v1.POST("/providers", routes.PostProvider)
		v1.POST("/products", routes.PostProduct)
		v1.POST("/accounts", routes.PostAccount)
		v1.POST("/tags", routes.PostTag)
		v1.POST("/sale_detail", routes.PostSaleDetail)
		v1.POST("/purchase_detail", routes.PostPurchaseDetail)
		v1.POST("/sales", routes.PostSale)
		v1.POST("/purchases", routes.PostPurchase)
		v1.POST("tags_customer", routes.PostTagCustomer)

		// *** Admin and manager ***
		// Stats
		v1.POST("/productsrank-k/:k", routes.GetRankProductK)
		v1.POST("/productsrank-cs/:k/:category", routes.GetRankProductCategoryS)
		v1.POST("/productsrank-cp/:k/:category", routes.GetRankProductCategoryP)
		v1.POST("/productsrank-b/:k/:brand", routes.GetRankProductBrand)
		v1.POST("/productsrank-pp/:id_product", routes.GetRankProductPP)
		v1.POST("/productsrank-r/:k", routes.GetRankProfitability)

		v1.POST("/customersrank-k/:k", routes.GetRankCustomerK)
		v1.POST("/customersrank-p/:k/:l", routes.GetRankCustomerKL)
		v1.POST("/customersrank-v/:k", routes.GetRankCustomerVariety)
		v1.POST("/customersrank-f/:k", routes.GetRankFrequency)

		v1.POST("/purchasesrank-k/:k", routes.GetRankPurchasesK)
		v1.POST("/purchasesrank-cp/:k/:category", routes.GetRankPurchasesCP)
		v1.POST("/purchasesrank-p/:k", routes.GetRankPurchasesProduct)

		v1.POST("/providersrank-k/:k", routes.GetRankProviderK)
		v1.POST("/providersrank-v/:k", routes.GetRankProviderVariety)
		v1.POST("/providersrank-pp/:k/:id_provider", routes.GetRankProviderPP)

		v1.POST("/salesrank-k/:k", routes.GetRankSalesK)
		v1.POST("/salesrank-c/:k/:category", routes.GetRankSalesCategory)
		v1.POST("/salesrank-p/:k", routes.GetRankSalesProduct)
		v1.POST("/salesrank-r/:k", routes.GetRankSalesArea)

		// Record
		v1.POST("/productsrec/:id", routes.GetSalesProductIDRec)

		v1.POST("/customersrec-p/:id_customer", routes.GetProductTotal)
		v1.POST("/customersrec-c/:id_customer", routes.GetTotalCash)

		v1.POST("/purchasesrec-p/:id_product", routes.GetPurchasesProduct)

		v1.POST("/salesrec-p/:id_product", routes.GetSalesProduct)
		v1.POST("/sales-total", routes.GetSales)

		// *** Seller ***
		// Stats
		v1.POST("/sellerproductsrank-k/:k/:seller")
		v1.POST("/sellerproductsrank-c/:category")
		v1.POST("/sellerproductsrank-b/:brand")

		v1.POST("/sellercustomer")
		v1.POST("/sellerc")
		v1.POST("/seller")

		v1.POST("/sellersalesrank-k/:k")
		v1.POST("/sellersalesrank-c/:k/:category")
		v1.POST("/sellersalesrank-p/:k")

		// Record
		v1.POST("/sales/:mail", routes.GetSalesID)

		// *** Dashboard ***
		v1.GET("/dashboard-info/:role/:id_seller", routes.GetInformationDashboard)
	}
	r.Run(":" + port)
}
