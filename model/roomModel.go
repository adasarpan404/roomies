package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	User       primitive.ObjectID   `bson:"user" validate:"required"`
	Title      *string              `bson:"title" validate:"required"`
	Address    *string              `bson:"address" validate:"required"`
	CreatedAt  time.Time            `bson:"createdAt"`
	Interested []primitive.ObjectID `bson:"interested,omitempty"`
	UpdatedAt  time.Time            `bson:"updatedAt"`
}
