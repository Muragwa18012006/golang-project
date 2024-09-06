package routes

import (
	"booking-app/controllers"
	"booking-app/middleware"

	"github.com/gin-gonic/gin"
)

func GetTodos(rt *gin.Engine) {
	protected := rt.Group("/users")
	protected.Use(middleware.Authenticate())
	protected.GET("/post/:id", controllers.GetAllTodos())
	protected.GET("/todos/:id", controllers.GetTodo())
	protected.DELETE("/:id", controllers.DeleteTodo())
}
