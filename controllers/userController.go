package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adasarpan404/roomies-be/helper"
	"github.com/adasarpan404/roomies-be/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, ok := c.Get("userId")
		if !ok {
			helper.ErrorResponse(c, http.StatusBadRequest, "User ID not found in context")

			return
		}

		objectUserId, err := primitive.ObjectIDFromHex(fmt.Sprint(userId))
		if err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user model.User
		projection := bson.M{"password": 0}
		err = userCollection.FindOne(ctx, bson.M{"_id": objectUserId}, options.FindOne().SetProjection(projection)).Decode(&user)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
