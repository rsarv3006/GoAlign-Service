package handler

import (
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeamSettings(dto model.TeamSettingsCreateDto) (*model.TeamSettings, error) {
	query := database.TeamSettingsCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	var teamSettings *model.TeamSettings

	rows, err := stmt.Query(dto.TeamId, dto.CanAllMembersAddTasks)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		teamSettings = &model.TeamSettings{}
		err := rows.Scan(
			&teamSettings.TeamSettingsId,
			&teamSettings.TeamId,
			&teamSettings.CanAllMembersAddTasks,
		)

		if err != nil {
			return nil, err
		}
	}

	return teamSettings, nil
}

func DeleteTeamSettingsByTeamId(teamId uuid.UUID) error {
	query := database.TeamSettingsDeleteByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(teamId)

	if err != nil {
		return err
	}

	return nil
}
