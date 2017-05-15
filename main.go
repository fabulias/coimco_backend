package main

import (
	"coimco_backend/routes"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	//"strings"
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
	println("Hola mundo el puerto es: " + port)
	// initialize the DbMap

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(Cors())
	// Simple group: v1
	v1 := r.Group("/v1")
	{
		v1.GET("/port", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"Port": port})
		})

		v1.GET("/customers", routes.GetCustomers)
		v1.GET("/providers", routes.GetProviders)

		v1.POST("/providers", routes.PostProvider)

	}
	r.Run(":" + port)
}
