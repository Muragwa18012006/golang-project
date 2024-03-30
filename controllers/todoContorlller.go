package controllers

import (
	"booking-app/database"
	"booking-app/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTodos() gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := database.Collection.Find(c, bson.M{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		/* var todos model.TODO */
		var allUsers []bson.M
		if err = cursor.All(c, &allUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, allUsers)
	}
}
func AddTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo model.CreateTodo
		if err := c.ShouldBind(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		result, err := database.Collection.InsertOne(c, todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "todo not created"})
			return
		}
		todos := model.TODO{
			ID:          result.InsertedID.(primitive.ObjectID),
			Title:       todo.Title,
			Description: todo.Description,
		}
		c.JSON(http.StatusOK, todos)
	}
}
func DeleteTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid path url"})
			return
		}
		res, err := database.Collection.DeleteOne(c, bson.M{"_id": _id})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "there was a problem trying to delete a todo "})
			return
		}
		if res.DeletedCount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "there was no such todo to delete"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": " a TODO is successfully deleted"})
	}
}
func UpdateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid path url"})
			return
		}
		var body struct {
			Title       string `json:"Title" bson:"Title" binding:"required"`
			Description string `json:"description" bson:"description" binding:"required"`
		}
		if err := c.ShouldBind(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid body request"})
			return
		}
		_, err = database.Collection.UpdateOne(c, bson.M{"_id": _id}, bson.M{"$set": bson.M{"Title": body.Title, "description": body.Description}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to update Todo"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Todo is successfully updated"})
	}
}
func GetTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID url"})
			return
		}
		result := database.Collection.FindOne(c, bson.M{"_id": _id})
		todo := model.TODO{}
		err = result.Decode(&todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to find a blog"})
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}
