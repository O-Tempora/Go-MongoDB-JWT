package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	GUID         primitive.ObjectID `bson:"_id"`
	Name         string             `json:"name" validate:"required,min=3"`
	RefreshToken string             `json:"refresh"`
}
