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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user ID from context
		userId, ok := c.Get("_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
			return
		}

		// Convert user ID to ObjectID
		objectUserId, err := primitive.ObjectIDFromHex(fmt.Sprint(userId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel() // Ensure cancel is called to release resources

		// Retrieve user from the database
		var user model.User
		projection := bson.M{"password": 0}
		err = userCollection.FindOne(ctx, bson.M{"_id": objectUserId}, options.FindOne().SetProjection(projection)).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return user in the response
		c.JSON(http.StatusOK, user)
	}
}
