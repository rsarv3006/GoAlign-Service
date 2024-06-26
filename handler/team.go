package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func GetTeamsForCurrentUser(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	query := database.TeamGetByUserIdQueryString

	rows, err := database.POOL.Query(context.Background(), query, currentUser.UserId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	teams := make([]model.Team, 0)
	teamIds := make([]uuid.UUID, 0)

	for rows.Next() {
		team := model.Team{}
		err := rows.Scan(
			&team.TeamId,
			&team.TeamName,
			&team.CreatorUserId,
			&team.Status,
			&team.TeamManagerId,
			&team.CreatedAt,
			&team.UpdatedAt,
		)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		teams = append(teams, team)
		teamIds = append(teamIds, team.TeamId)
	}

	users, err := getUsersByTeamIdArray(teamIds)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	tasks, err := getTasksByTeamIdArray(teamIds)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamReturnArray := make([]model.TeamReturnWithUsersAndTasks, 0)

	for _, team := range teams {
		teamValue := team

		teamTasks, tasksOk := tasks[team.TeamId]

		if !tasksOk {
			teamTasks = make([]model.TaskReturnWithTaskEntries, 0)
		}

		teamUsers, usersOk := users[team.TeamId]

		if !usersOk {
			teamUsers = make([]model.User, 0)
		}

		team := model.TeamReturnWithUsersAndTasks{
			Team:  &teamValue,
			Users: teamUsers,
			Tasks: teamTasks,
		}

		teamReturnArray = append(teamReturnArray, team)
	}

	defer func() {
		teamReturnArray = make([]model.TeamReturnWithUsersAndTasks, 0)
		teams = make([]model.Team, 0)
		teamIds = make([]uuid.UUID, 0)
		users = make(map[uuid.UUID][]model.User)
		tasks = make(map[uuid.UUID][]model.TaskReturnWithTaskEntries)
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"teams":   teamReturnArray,
	})
}

func CreateTeam(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	teamDto := new(model.TeamCreateDto)

	if err := c.BodyParser(teamDto); err != nil {
		return sendBadRequestResponse(c, err, "Unable to parse team create dto")
	}

	if teamDto.TeamName == "" {
		err := errors.New("Team name is required")
		return sendBadRequestResponse(c, err, "Team name is required")
	}

	query := database.TeamCreateQueryString

	team := new(model.Team)

	cleanedTeamName := helper.SanitizeInput(teamDto.TeamName)
	rows, err := database.POOL.Query(context.Background(), query, cleanedTeamName, currentUser.UserId, currentUser.UserId)

	if err != nil {
		switch e := err.(type) {
		case *pq.Error:
			fmt.Println("Postgres error:", e.Message)

		default:
			fmt.Println("Unknown error:", e)

		}
		return sendBadRequestResponse(c, err, "")
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&team.TeamId, &team.TeamName, &team.CreatorUserId, &team.Status, &team.TeamManagerId, &team.CreatedAt, &team.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	_, err = CreateUserTeamMembership(currentUser.UserId, team.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamUsers := make([]model.User, 0)
	teamUsers = append(teamUsers, *currentUser)

	teamTasks := make([]model.TaskReturnWithTaskEntries, 0)

	teamReturn := model.TeamReturnWithUsersAndTasks{
		Team:  team,
		Users: teamUsers,
		Tasks: teamTasks,
	}

	teamSettingsDto := new(model.TeamSettingsCreateDto)
	teamSettingsDto.TeamId = team.TeamId
	teamSettingsDto.CanAllMembersAddTasks = false
	teamSettings, err := CreateTeamSettings(*teamSettingsDto)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer func() {
		teamDto = nil
		team = nil
		teamUsers = nil
		teamTasks = nil
		teamReturn = model.TeamReturnWithUsersAndTasks{}
		teamSettingsDto = nil
		teamSettings = nil
	}()

	return c.Status(201).JSON(&fiber.Map{
		"message":  "Team created successfully",
		"team":     teamReturn,
		"settings": teamSettings,
	})
}

func getTeamsByTeamIdArray(teamIds []uuid.UUID) (map[uuid.UUID]model.TeamReturnWithUsersAndTasks, error) {
	query := database.TeamGetByIdsQueryString

	rows, err := database.POOL.Query(context.Background(), query, pq.Array(teamIds))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	teams := make([]model.Team, 0)

	for rows.Next() {
		team := model.Team{}
		err := rows.Scan(&team.TeamId, &team.TeamName, &team.CreatorUserId, &team.Status, &team.TeamManagerId, &team.CreatedAt, &team.UpdatedAt)

		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	tasks, err := getTasksByTeamIdArray(teamIds)

	if err != nil {
		return nil, err
	}

	users, err := getUsersByTeamIdArray(teamIds)

	if err != nil {
		return nil, err
	}

	teamReturnMap := make(map[uuid.UUID]model.TeamReturnWithUsersAndTasks)

	for _, team := range teams {
		teamUsers, usersOk := users[team.TeamId]

		if !usersOk {
			teamUsers = make([]model.User, 0)
		}

		teamTasks, tasksOk := tasks[team.TeamId]

		if !tasksOk {
			teamTasks = make([]model.TaskReturnWithTaskEntries, 0)
		}

		teamReturn := model.TeamReturnWithUsersAndTasks{
			Team:  &team,
			Users: teamUsers,
			Tasks: teamTasks,
		}

		teamReturnMap[team.TeamId] = teamReturn
	}

	return teamReturnMap, nil
}

func getTeamById(teamId uuid.UUID) (*model.TeamReturnWithUsersAndTasks, error) {
	teams, err := getTeamsByTeamIdArray([]uuid.UUID{teamId})

	if err != nil {
		return nil, err
	}

	team, ok := teams[teamId]

	if !ok {
		return nil, errors.New("Team not found")
	}

	return &team, nil
}

func DeleteTeam(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	teamId, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	teamToDelete, err := getTeamById(teamId)

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	if teamToDelete.TeamManagerId != currentUser.UserId {
		return sendForbiddenResponse(c)
	}

	err = DeleteTeamSettingsByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteUserTeamMembershipsByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteTaskEntriesByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteTasksByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteTeamInvitesByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteTeam(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func deleteTeam(teamId uuid.UUID) error {

	_, err := database.POOL.Exec(context.Background(), database.TeamDeleteQueryString, teamId)

	if err != nil {
		return err
	}

	return nil
}

func isUserATeamManagerOfAnyTeam(userId uuid.UUID) bool {
	query := database.TeamGetByTeamManagerIdQueryString

	rows, err := database.POOL.Query(context.Background(), query, userId)

	if err != nil {
		return false
	}

	defer rows.Close()

	for rows.Next() {
		return true
	}

	return false
}

func GetTeamByTeamIdEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	teamId, err := uuid.Parse(c.Params("teamId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	team, err := getTeamById(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	isUserInTeam, err := isUserInTeam(currentUser.UserId, teamId)

	if err != nil {
		return sendBadRequestResponse(c, err, "User is not in team")
	}

	if !isUserInTeam {
		return sendForbiddenResponse(c)
	}

	return c.Status(200).JSON(&fiber.Map{
		"team": team,
	})
}

func UpdateTeamManagerEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

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

	teamManagerId, err := uuid.Parse(c.Params("teamManagerId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team manager id")
	}

	isUserInTeam, err := isUserInTeam(teamManagerId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserInTeam {
		err = errors.New("User is not in team")
		return sendBadRequestResponse(c, err, "User is not in team")
	}

	query := database.TeamUpdateTeamManagerQueryString

	rows, err := database.POOL.Query(context.Background(), query, teamManagerId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	return c.Status(200).JSON(&fiber.Map{
		"message": "Team manager updated successfully",
	})
}
