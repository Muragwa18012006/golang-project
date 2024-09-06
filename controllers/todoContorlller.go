package controllers

import (
	"booking-app/database"
	"booking-app/helper"
	"booking-app/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func GetAllTodos() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* cursor, err := database.Collection.Find(c, bson.M{})  */
		id := c.Param("id")
		_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		alls := bson.M{"UserID": _id}
		cursor, err := database.Collection.Find(c, alls)

		/*	todo := model.TODO{}
			err = result.Decode(&todo) */
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		/* c.JSON(http.StatusOK, todo) */
		/* if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		} */
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
			UserID:      todo.UserId,
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

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		token, refreshToken := helper.GenerateAllTokens(user.Firstname, user.LastName, user.Email, *user.Password)
		user.Token = &token
		user.RefreshToken = &refreshToken
		result, err := database.UserCollection.InsertOne(c, user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "user not created"})
			return
		}
		userID := result.InsertedID.(primitive.ObjectID)
		c.JSON(http.StatusOK, userID)
	}
}
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		var foundUser model.User

		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := database.UserCollection.FindOne(c, bson.M{"Email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if foundUser.Email == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		token, refreshToken := helper.GenerateAllTokens(foundUser.Firstname, foundUser.LastName, foundUser.Email, *foundUser.Password)
		helper.UpdateAllTokens(token, refreshToken, *foundUser.Password)
		/* err := database.UserCollection.FindOne(c, bson.M{"Password":foundUser.Password})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		} */
		c.JSON(http.StatusOK, foundUser)
	}
}
