package auth

import (
	"github.com/dgrijalva/jwt-go/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtWrapper struct {
	SecretKey      string
	Issuer         string
	ExpirationTime int64
}

type JwtClaim struct {
	UserId   primitive.ObjectID
	Email    string
	UserType string
	jwt.StandardClaims
}
