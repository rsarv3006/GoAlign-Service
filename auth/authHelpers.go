package auth

import (
	"context"
	"errors"
	"time"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

type JWTClaims struct {
	User model.User
	jwt.StandardClaims
}

const SevenDays = 7 * 24 * time.Hour

var (
	ErrExpired = errors.New("token expired")
	ErrInvalid = errors.New("couldn't parse claims")
)

func GenerateJWT(user model.User, ctx *fiber.Ctx) (*string, error) {
	jwtSecretString := ctx.Locals("JwtSecret").(string)
	jwtKey := []byte(jwtSecretString)

	expirationTime := time.Now().Add(SevenDays)
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

func ValidateToken(signedToken string, ctx *fiber.Ctx) (*model.User, error) {
	jwtSecretString := ctx.Locals("JwtSecret").(string)
	jwtKey := []byte(jwtSecretString)

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
		return nil, ErrInvalid
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, ErrExpired
	}
	return &claims.User, nil
}

func SendEmailWithCode(ctx *fiber.Ctx, code string, username string, email string) (bool, error) {
	var vanillaContext context.Context
	br := ctx.Locals("BrevoClient").(*brevo.APIClient)

	_, _, err := br.TransactionalEmailsApi.SendTransacEmail(vanillaContext, brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "GoAlign",
			Email: "goalign.app@gmail.com",
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Name:  username,
				Email: email,
			},
		},
		Subject:     "GoAlign Authorization Code",
		TextContent: CreateNonHtmlStringForUserAuth(username, code),
		HtmlContent: CreateEmailForUserAuth(username, code),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
