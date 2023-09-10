package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateUserTeamMembership(userId uuid.UUID, teamId uuid.UUID) (*model.UserTeamMembership, error) {
	query := database.UserTeamMembershipCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	userTeamMembership := new(model.UserTeamMembership)

	rows, err := stmt.Query(userId, teamId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&userTeamMembership.UserTeamMembershipId, &userTeamMembership.UserId, &userTeamMembership.TeamId, &userTeamMembership.CreatedAt, &userTeamMembership.UpdatedAt, &userTeamMembership.Status)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return userTeamMembership, nil
}

func getUserTeamMemberships(teamId uuid.UUID) ([]model.UserTeamMembership, error) {
	query := database.UserTeamMembershipGetByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(teamId)

	if err != nil {
		return nil, err
	}

	var users []model.UserTeamMembership

	for rows.Next() {
		userTeamMembership := new(model.UserTeamMembership)
		err := rows.Scan(&userTeamMembership.UserTeamMembershipId, &userTeamMembership.UserId, &userTeamMembership.TeamId, &userTeamMembership.CreatedAt, &userTeamMembership.UpdatedAt, &userTeamMembership.Status)

		if err != nil {
			return nil, err
		}

		users = append(users, *userTeamMembership)
	}

	return users, nil

}

func isUserInTeam(userId uuid.UUID, teamId uuid.UUID) (bool, error) {
	query := database.UserTeamMembershipGetByUserIdAndTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userId, teamId)

	if err != nil {
		return false, err
	}

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func deleteUserTeamMembershipsByTeamId(teamId uuid.UUID) error {
	query := database.UserTeamMembershipDeleteByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(teamId)

	if err != nil {
		return err
	}

	return nil
}

func deleteUserTeamMembershipsByUserId(userId uuid.UUID) error {
	query := database.UserTeamMembershipDeleteByUserIdQuery
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

func RemoveUserFromTeamEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	teamId, err := uuid.Parse(c.Params("teamId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	team, err := getTeamById(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if team.TeamManagerId != currentUser.UserId {
		return sendForbiddenResponse(c)
	}

	userId, err := uuid.Parse(c.Params("userId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid user id")
	}

	isUserInTeam, err := isUserInTeam(userId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserInTeam {
		err = errors.New("User is not in team")
		return sendBadRequestResponse(c, err, "User is not in team")
	}

	if team.TeamManagerId == userId {
		err = errors.New("Cannot remove team manager from team")
		return sendBadRequestResponse(c, err, "Cannot remove team manager from team")
	}

	query := database.UserTeamMembershipDeleteQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	_, err = stmt.Query(userId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func getTeamMembersByTeamId(teamId uuid.UUID) ([]model.User, error) {
	// userTeamMemberships, err := getUserTeamMemberships(teamId)

	// if err != nil {
	// return nil, err
	// }

	// userIds := make([]uuid.UUID, len(userTeamMemberships))

	// return users, nil

	return nil, nil
}
