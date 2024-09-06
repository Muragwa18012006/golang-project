package main

import (
	"booking-app/database"
	"booking-app/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	router *gin.Engine
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
	/* router.Use(cors.Default()) */
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000","https://next-js-blog-phi-ruddy.vercel.app"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "token"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * 60 * 60,
	}))
	router.Use(gin.Logger())
	routes.AddTodos(router)
	routes.GetTodos(router)
	routes.AuthRoutes(router)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "route testing ...... successfully"})
	})
	router.Run("0.0.0.0:" + port)
}
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
