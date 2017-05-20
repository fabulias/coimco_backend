package main

import (
	"coimco_backend/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func Cors() gin.HandlerFunc {
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
	r.GET("/login", routes.Login)
	// Simple group: v1
	v1 := r.Group("api")
	{
		//Methods plurals GET
		v1.GET("/customers", routes.GetCustomers)
		v1.GET("/products", routes.GetProducts)
		//v1.GET("/providers", routes.GetProviders)

		//Methods singular GET
		v1.GET("/customers/:mail", routes.GetCustomer)
		v1.GET("products/:id", routes.GetProduct)

		//Methods POST
		v1.POST("/customers", routes.PostCustomer)
		v1.POST("/products", routes.PostProduct)

	}
	r.Run(":" + port)
}
