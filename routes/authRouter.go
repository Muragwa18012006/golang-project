package routes

import (
	"booking-app/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rt *gin.Engine) {
	rt.POST("/create", controllers.CreateUser())
	rt.POST("/login", controllers.Login())
}
