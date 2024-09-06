package routes

import (
	"booking-app/controllers"
	"booking-app/middleware"

	"github.com/gin-gonic/gin"
)

func AddTodos(rt *gin.Engine) {
	protected := rt.Group("/users")
	protected.Use(middleware.Authenticate())
	protected.POST("/addTodos", controllers.AddTodo())
	protected.PUT("/update/:id", controllers.UpdateTodo())

}
