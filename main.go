package main

import (
	"booking-app/database"
	"booking-app/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := database.DbConnect()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connected to mongodb")
	defer func() {
		err := database.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("error occured while loading env")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	router := gin.New()
	router.Use(cors.Default())
	router.Use(gin.Logger())
	routes.AddTodos(router)
	routes.GetTodos(router)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "route testing ...... successfully"})
	})
	router.Run(":" + port)
}
