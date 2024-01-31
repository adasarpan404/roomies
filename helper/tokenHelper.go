package helper

import (
	// "fmt"
	// "log"
	"time"

	"github.com/adasarpan404/roomies-be/environment"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	ID        string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, firstName string, lastName string, id string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		ID:        id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(environment.SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
