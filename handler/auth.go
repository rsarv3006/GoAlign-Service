package handler

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

// TODO: add check to make sure we don't create duplicate login requests

func Register(c *fiber.Ctx) error {
	userCreateDto := new(model.UserCreateDto)

	if err := c.BodyParser(userCreateDto); err != nil {
		log.Println(err)
		return sendBadRequestResponse(c, err, "Invalid request body")
	}

	user := new(model.User)

	userName := helper.SanitizeInput(userCreateDto.UserName)

	user.UserName = userName
	user.Email = userCreateDto.Email

	if user.UserName == "" || user.Email == "" {
		err := fmt.Errorf("Username and Email are required")
		return sendBadRequestResponse(c, err, "Username and Email are required")
	}

	user.UserId = uuid.New()
	user.IsActive = true

	// TODO: conver to prepared statement
	rows, err := database.DB.Query("INSERT INTO users (username, email ) VALUES ($1, $2) RETURNING *",
		user.UserName,
		user.Email,
	)

	if err != nil {
		if strings.Contains(err.Error(), `"users_email_key"`) {
			return sendBadRequestResponse(c, err, "Email already exists")
		}
		return sendInternalServerErrorResponse(c, err)
	}

	if rows.Next() {
		err := rows.Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified, &user.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return sendInternalServerErrorResponse(c, err)
		}
	}

	loginRequest, err := createLoginRequest(user.UserId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "login_request_id": loginRequest.LoginRequestId, "user_id": loginRequest.UserId})
}

func Login(c *fiber.Ctx) error {
	loginInitiateDto := new(model.LoginInitiateDto)

	if err := c.BodyParser(loginInitiateDto); err != nil {
		log.Println(err)
		return sendBadRequestResponse(c, err, "Invalid request body")
	}

	query := database.UserGetUserByEmailQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	user := model.User{}

	err = stmt.QueryRow(loginInitiateDto.Email).Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return sendNotFoundResponse(c, "User not found")
		}
		return sendInternalServerErrorResponse(c, err)
	}

	loginRequest, err := createLoginRequest(user.UserId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "login_request_id": loginRequest.LoginRequestId, "user_id": loginRequest.UserId})
}

func FetchCode(c *fiber.Ctx) error {
	loginRequestDto := new(model.LoginCodeRequestDto)

	if err := c.BodyParser(loginRequestDto); err != nil {
		return sendBadRequestResponse(c, err, "Invalid request body")
	}

	if loginRequestDto.LoginRequestToken == "" {
		err := fmt.Errorf("Code is required")
		return sendBadRequestResponse(c, err, "Code is required")
	}

	loginRequestQuery := database.LoginRequestGetByLoginRequestId
	stmt, err := database.DB.Prepare(loginRequestQuery)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	loginRequest := model.LoginRequest{}

	err = stmt.QueryRow(loginRequestDto.LoginCodeRequestId).Scan(
		&loginRequest.LoginRequestId,
		&loginRequest.UserId,
		&loginRequest.LoginRequestDate,
		&loginRequest.LoginRequestExpirationDate,
		&loginRequest.LoginRequestToken,
		&loginRequest.LoginRequestStatus)

	if err != nil {
		log.Println(err)
		log.Println("confirm location")
		return sendUnauthorizedResponse(c)
	}

	if loginRequest.LoginRequestStatus != "pending" {
		log.Println(err)
		return sendUnauthorizedResponse(c)
	}

	if loginRequest.LoginRequestExpirationDate.Before(time.Now()) {
		log.Println(err)
		return sendUnauthorizedResponse(c)
	}

	if loginRequest.LoginRequestToken != loginRequestDto.LoginRequestToken {
		log.Println(err)
		return sendUnauthorizedResponse(c)
	}

	if loginRequest.UserId != loginRequestDto.UserId {
		log.Println(err)
		return sendUnauthorizedResponse(c)
	}

	userFromDb, errFromDb := database.DB.Query("SELECT * FROM users WHERE user_id = $1", loginRequest.UserId)

	if errFromDb != nil {
		log.Println(errFromDb)
		return sendUnauthorizedResponse(c)
	}

	defer userFromDb.Close()

	user := model.User{}

	for userFromDb.Next() {
		switch err := userFromDb.Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified, &user.CreatedAt); err {
		case sql.ErrNoRows:
			fmt.Println("sql.ErrNoRows")
			return sendUnauthorizedResponse(c)
		case nil:
			// Expected outcome, user found
		default:
			fmt.Println(err)
			return sendUnauthorizedResponse(c)
		}
	}

	if !isUserObjectNotNil(&user) {
		log.Println(err)
		return sendUnauthorizedResponse(c)
	}

	signedTokenString, err := auth.GenerateJWT(user)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = markLoginRequestAsCompleted(loginRequest.LoginRequestId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": signedTokenString})
}

func isUserObjectNotNil(user *model.User) bool {
	return user.UserId != uuid.Nil

}

func createLoginRequest(userId uuid.UUID) (*model.LoginRequest, error) {
	// TODO: Send email with code
	query := database.CreateLoginRequestQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	loginRequestExpirationDate := time.Now().Add(time.Minute * 10)
	loginCode, err := generateUniqueLoginCode()

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userId, loginRequestExpirationDate, loginCode)

	if err != nil {
		return nil, err
	}

	loginRequest := model.LoginRequest{}
	if rows.Next() {
		err := rows.Scan(&loginRequest.LoginRequestId, &loginRequest.UserId, &loginRequest.LoginRequestDate, &loginRequest.LoginRequestExpirationDate, &loginRequest.LoginRequestToken, &loginRequest.LoginRequestStatus)
		if err != nil {
			return nil, err
		}
	}

	return &loginRequest, nil
}

func generateUniqueLoginCode() (string, error) {
	const maxAttempts = 10
	var attempts int = 0

	for {
		if attempts >= maxAttempts {
			return "", fmt.Errorf("Could not generate unique login code after %d attempts", maxAttempts)
		}

		attempts++

		loginCode := helper.GenerateCodeHelper()
		rows, err := database.DB.Query("SELECT * FROM login_requests WHERE login_request_token = $1", loginCode)

		if err != nil {
			return "", err
		}

		if !rows.Next() {
			return loginCode, nil
		}
	}
}

func markLoginRequestAsCompleted(loginRequestId uuid.UUID) error {
	query := database.LoginRequestMarkAsCompletedQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(loginRequestId)

	if err != nil {
		return err
	}

	return nil
}
