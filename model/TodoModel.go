package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TODO struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"Title" bson:"Title"`
	Description string             `json:"description" bson:"description"`
}
type CreateTodo struct {
	Title       string `json:"Title" bson:"Title"`
	Description string `json:"description" bson:"description"`
}
