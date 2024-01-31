package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adasarpan404/roomies-be/db"
	"github.com/adasarpan404/roomies-be/helper"
	"github.com/adasarpan404/roomies-be/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = db.OpenCollection(db.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.TODO(), 100*time.Second)
		var user model.User
		if err := c.BindJSON(&user); err != nil {
			fmt.Print(user)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email alreay exists"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		token, err := helper.GenerateToken(*user.Email, *user.FirstName, *user.LastName, user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userObj, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			msg := "User item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, gin.H{"user": userObj, "token": token})
	}
}
