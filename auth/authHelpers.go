package auth

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

var jwtKey = []byte("supersecretkey")

type JWTClaims struct {
	User model.User
	jwt.StandardClaims
}

func GenerateJWT(user model.User) (*string, error) {
	log.Println(user)
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return &tokenString, err
}

func ValidateToken(signedToken string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}
	return &claims.User, nil
}
