package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TODO struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"Title" bson:"Title"`
	Description string             `json:"description" bson:"description"`
	UserID      primitive.ObjectID `json:"UserID" bson:"UserID"`
}
type CreateTodo struct {
	Title       string             `json:"Title" bson:"Title"`
	Description string             `json:"description" bson:"description"`
	UserId      primitive.ObjectID `json:"UserID" bson:"UserID"`
}
type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Firstname    string             `json:"Firstname" bson:"Firstname"`
	LastName     string             `json:"Lastname" bson:"Lastname"`
	Email        string             `json:"Email" bson:"Email"`
	Password     *string            `json:"Password" bson:"Password"`
	Token        *string            `json:"token" bson:"token"`
	RefreshToken *string            `json:"refreshToken" bson:"refreshToken"`
}
