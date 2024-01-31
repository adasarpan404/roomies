package helper

import (
	"log"
	"time"

	"github.com/adasarpan404/roomies-be/environment"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, firstName string, lastName string, id string) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["email"] = email
	claims["firstName"] = firstName
	claims["lastName"] = lastName
	claims["id"] = id
	tokenString, err := token.SignedString(environment.SECRET_KEY)
	if err != nil {
		log.Panic(err)
	}
	return tokenString, nil
}
