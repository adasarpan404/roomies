package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/adasarpan404/roomies-be/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user model.User
		err := userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
