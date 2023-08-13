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
