package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	defer rows.Close()

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

	defer rows.Close()

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

	defer rows.Close()

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

func isAllowedToDeleteUserFromTeam(currentUser *model.User, userId uuid.UUID, teamManagerId uuid.UUID) bool {
	if teamManagerId == currentUser.UserId {
		return true
	}

	if currentUser.UserId == userId {
		return true
	}

	return false
}

func RemoveUserFromTeamEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	teamId, err := uuid.Parse(c.Params("teamId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	team, err := getTeamById(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
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

	isAllowedToDeleteTeamInvite := isAllowedToDeleteUserFromTeam(currentUser, userId, team.TeamManagerId)

	if !isAllowedToDeleteTeamInvite {
		err = errors.New("Not allowed to delete user from team")
		return sendBadRequestResponse(c, err, "Not allowed to delete user from team")
	}

	query := database.UserTeamMembershipDeleteQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(userId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	return c.SendStatus(fiber.StatusNoContent)
}
