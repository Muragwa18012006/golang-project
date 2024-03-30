package routes

import (
	"booking-app/controllers"

	"github.com/gin-gonic/gin"
)

func GetTodos(rt *gin.Engine) {
	rt.GET("users/todos", controllers.GetAllTodos())
	rt.GET("users/todos/:id", controllers.GetTodo())
	rt.DELETE("users/:id", controllers.DeleteTodo())
}
