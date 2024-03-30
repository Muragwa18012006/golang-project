package routes

import (
	"booking-app/controllers"

	"github.com/gin-gonic/gin"
)

func AddTodos(rt *gin.Engine) {
	rt.POST("users/addTodos", controllers.AddTodo())
	rt.PUT("users/update/:id", controllers.UpdateTodo())
}
