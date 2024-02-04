package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/adasarpan404/roomies-be/helper"
	"github.com/adasarpan404/roomies-be/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		perPage, err := strconv.Atoi(c.Query("limit"))
		if err != nil || perPage < 1 {
			perPage = 10
		}

		skip := (page - 1) * perPage
		var ctx, cancel = context.WithTimeout(context.TODO(), 100*time.Second)
		totalCount, err := roomCollection.CountDocuments(ctx, bson.M{})
		defer cancel()
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		cursor, err := roomCollection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(perPage)))
		defer cancel()
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		var results []model.Room
		if err := cursor.All(ctx, &results); err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		hasPrevPage := page > 1

		// Check if there is a next page
		hasNextPage := (page-1)*perPage+len(results) < int(totalCount)
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":      true,
				"data":        results,
				"total":       totalCount,
				"hasPrevPage": hasPrevPage,
				"hasNextPage": hasNextPage,
			})
	}
}

func UpdateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.TODO(), 100*time.Second)
		var room model.Room
		if err := c.BindJSON(&room); err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		validationErr := validate.Struct(room)
		if validationErr != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, validationErr.Error())
		}
		userId, ok := c.Get("userId")
		roomId := c.Param("id")
		defer cancel()
		if !ok {
			helper.ErrorResponse(c, http.StatusBadRequest, "User ID not found in the context")
			return
		}
		objectUserId, err := primitive.ObjectIDFromHex(fmt.Sprint(userId))
		defer cancel()
		if err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID Format")
			return
		}
		objectRoomId, err := primitive.ObjectIDFromHex(roomId)
		defer cancel()
		if err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, "Invalid Room ID Format")
			return
		}

		err = roomCollection.FindOne(ctx, bson.M{"_id": roomId, "user": objectUserId}).Decode(&room)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		filter := bson.M{"_id": objectRoomId}
		update := bson.M{
			"$set": bson.M{
				"title":   room.Title,
				"address": room.Address,
			},
		}
		result, err := roomCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, "Error while Updating the document")
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.TODO(), 100*time.Second)
		var room model.Room
		if err := c.BindJSON(&room); err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		validationErr := validate.Struct(room)
		if validationErr != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, validationErr.Error())
		}

		userId, ok := c.Get("userId")
		defer cancel()
		if !ok {
			helper.ErrorResponse(c, http.StatusBadRequest, "User ID not found in the context")
			return
		}

		objectUserId, err := primitive.ObjectIDFromHex(fmt.Sprint(userId))
		defer cancel()
		if err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
			return
		}
		room.User = objectUserId
		room.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		room.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		room.ID = primitive.NewObjectID()
		roomObj, err := roomCollection.InsertOne(ctx, room)
		if err != nil {
			msg := "Room item was not created"
			helper.ErrorResponse(c, http.StatusInternalServerError, msg)
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
