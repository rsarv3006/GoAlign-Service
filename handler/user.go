package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func DeleteUserEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	userId := currentUser.UserId

	err := deleteTeamInvitesByUserEmail(currentUser.Email)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	isUserATeamManagerOfAnyTeam := isUserATeamManagerOfAnyTeam(userId)

	if isUserATeamManagerOfAnyTeam {

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "User is a team manager of at least one team",
		})
	}

	err = deleteUserTeamMembershipsByUserId(userId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteUserLoginRequestsByUserId(userId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	query := database.UserDeleteUserQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(userId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func getUsersByTeamId(teamId uuid.UUID) ([]model.User, error) {
	query := database.UserGetUsersByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	var users []model.User

	if err != nil {
		return users, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(teamId)

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified, &user.CreatedAt)

		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func getUsersByIdArray(userIds []uuid.UUID) (map[uuid.UUID]model.User, error) {
	users := make(map[uuid.UUID]model.User)

	rows, err := database.POOL.Query(
		context.Background(),
		database.UserGetUsersByIdArrayQuery,
		pq.Array(userIds))

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified, &user.CreatedAt)

		if err != nil {
			return users, err
		}

		users[user.UserId] = user
	}

	return users, nil
}

func getUsersByTeamIdArray(teamIds []uuid.UUID) (map[uuid.UUID][]model.User, error) {
	query := database.UserGetUsersByTeamIdArrayQuery
	stmt, err := database.DB.Prepare(query)

	users := make(map[uuid.UUID][]model.User)

	if err != nil {
		return users, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(pq.Array(teamIds))

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User
		var teamId uuid.UUID
		err := rows.Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified, &user.CreatedAt, &teamId)

		if err != nil {
			return users, err
		}

		users[teamId] = append(users[teamId], user)
	}

	return users, nil
}

func getUserById(userId uuid.UUID) (model.User, error) {
	query := database.UserGetUserByIdQuery
	stmt, err := database.DB.Prepare(query)

	var user model.User

	if err != nil {
		return user, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(userId).Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified, &user.CreatedAt)

	if err != nil {
		return user, err
	}

	return user, nil
}

func deleteUserLoginRequestsByUserId(userId uuid.UUID) error {
	query := database.UserDeleteUserLoginRequestsByUserIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(userId)

	if err != nil {
		return err
	}

	return nil
}
