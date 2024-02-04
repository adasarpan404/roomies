package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adasarpan404/roomies-be/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.TODO(), 100*time.Second)
		cursor, err := roomCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var results []model.Room
		if err := cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": results})
	}
}

func UpdateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.TODO(), 100*time.Second)
		var room model.Room
		if err := c.BindJSON(&room); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		validationErr := validate.Struct(room)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		userId, ok := c.Get("userId")
		defer cancel()
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in the context"})
			return
		}

		objectUserId, err := primitive.ObjectIDFromHex(fmt.Sprint(userId))
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}
		room.User = objectUserId
		room.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		room.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		room.ID = primitive.NewObjectID()
		roomObj, err := roomCollection.InsertOne(ctx, room)
		if err != nil {
			msg := "Room item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, gin.H{"room": roomObj})
	}
}

func DeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
