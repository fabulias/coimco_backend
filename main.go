package main

import (
	"log"
	"os"

	"coimco_backend/auth"
	"coimco_backend/routes"
	"github.com/gin-gonic/gin"
)

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func Cors() gin.HandlerFunc {
	log.Println("CORS Middleware")
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
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
		v1.GET("/products/:id", routes.GetProduct)
		v1.GET("/accounts/:mail", routes.GetAccount)
		v1.GET("/tags/:id", routes.GetTag)
		v1.GET("/sale_detail/:sale_id/:product_id", routes.GetSaleDetail)
		v1.GET("/purchase_detail/:purchase_id/:product_id", routes.GetPurchaseDetail)
		v1.GET("/sales", routes.GetSale)
		v1.GET("/purchases", routes.PostPurchase)

		//Methods POST
		v1.POST("/customers", routes.PostCustomer)
		v1.POST("/products", routes.PostProduct)
		v1.POST("/accounts", routes.PostAccount)
		v1.POST("/tags", routes.PostTag)
		v1.POST("/sale_detail", routes.PostSaleDetail)
		v1.POST("/purchase_detail", routes.PostPurchaseDetail)
		v1.POST("/sales", routes.PostSale)
		v1.POST("/purchases", routes.PostPurchase)
	}
	r.Run(":" + port)
}
