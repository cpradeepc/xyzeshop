package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// generate token
func (j *JwtWrapper) GenerateToken(id primitive.ObjectID, email string, userType string) (token string, err error) {
	t := time.Now().Add(time.Hour * time.Duration(j.ExpirationTime))
	claims := &JwtClaim{
		UserId:   id,
		UserType: userType,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{Time: t},
			Issuer:    j.Issuer,
		},
	}

	tokenStr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenStr.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// validate token
func (j JwtWrapper) ValidateToken(signedToken string) (claim *JwtClaim, err error) {
	token, er := jwt.ParseWithClaims(
		signedToken, &JwtClaim{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if er != nil {
		log.Println("error in parse token : ", er)
		return nil, er
	}

	claims, ok := token.Claims.(*JwtClaim) //if assert succeed return true
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	log.Println("exp time ", claims.ExpiresAt.Unix())
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}
	return claims, nil
}

// check header's contains authorization or not, validate the token
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("you are not authorized"))
			return
		}

		extractedToken := strings.Split(token, "Bearer ")
		if len(extractedToken) != 2 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("you are not authorized"))
			return
		}

		token = strings.TrimSpace(extractedToken[1]) //token string without Bearer
		jwt_wrapper := JwtWrapper{
			SecretKey: os.Getenv("JwtScrets"),
			Issuer:    os.Getenv("JwtIssuer"),
		}

		claims, err := jwt_wrapper.ValidateToken(token)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("you are not authorized"))
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Next() //?

	}
}
