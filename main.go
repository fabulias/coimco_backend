package main

import (
	"coimco_backend/routes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Ping": "Pong"})
	})
	// Simple group: v1
	v1 := r.Group("/v1")
	{
		v1.GET("/customers", routes.GetCustomers)
		v1.GET("/products", routes.GetProducts)
		v1.GET("/providers", routes.GetProviders)

		v1.GET("/customers/:mail", routes.GetCustomer)

		v1.POST("/customers", routes.PostCustomers)

	}
	r.Run(":" + port)
}
