package handler

import (
	"fmt"
	"log"

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

	log.Println(userId)
	log.Println(teamId)

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
