package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func createTaskEntry(taskEntryCreateDto *model.TaskEntryCreateDto, currentUserId uuid.UUID) (*model.TaskEntryReturnWithAssignedUser, error) {
	task, err := getTaskByTaskId(taskEntryCreateDto.TaskId)

	if err != nil {
		return nil, err
	}

	isUserInTeam, err := isUserInTeam(taskEntryCreateDto.AssignedUserId, task.TeamId)

	if err != nil {
		return nil, err
	}

	if !isUserInTeam {
		return nil, errors.New("User is not in the team")
	}
	isUserTeamManager, err := isUserTheTeamManager(currentUserId, task.TeamId)

	if err != nil {
		return nil, err
	}

	teamSettings, err := getTeamSettingsByTeamId(task.TeamId)

	if err != nil {
		return nil, err
	}

	if !teamSettings.CanAllMembersAddTasks && !isUserTeamManager {
		return nil, errors.New("User is not the team manager")
	}

	query := database.TaskEntryCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	taskEntry := new(model.TaskEntry)

	rows, err := stmt.Query(
		taskEntryCreateDto.StartDate,
		taskEntryCreateDto.EndDate,
		taskEntryCreateDto.Notes,
		taskEntryCreateDto.AssignedUserId,
		taskEntryCreateDto.TaskId,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return nil, err
		}
	}

	assignedUser, err := getUserById(taskEntry.AssignedUserId)

	if err != nil {
		return nil, err
	}

	taskEntryReturnWithAssignedUser := model.TaskEntryReturnWithAssignedUser{
		TaskEntry:    taskEntry,
		AssignedUser: assignedUser,
	}

	return &taskEntryReturnWithAssignedUser, nil
}

func getTaskEntriesByTeamId(teamId uuid.UUID) ([]model.TaskEntryReturnWithAssignedUser, error) {
	query := database.TaskEntryGetByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	taskEntries := make([]model.TaskEntryReturnWithAssignedUser, 0)

	rows, err := stmt.Query(teamId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		taskEntry := new(model.TaskEntry)

		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return nil, err
		}

		assignedUser, err := getUserById(taskEntry.AssignedUserId)

		if err != nil {
			return nil, err
		}

		taskEntryReturnWithCreator := model.TaskEntryReturnWithAssignedUser{
			TaskEntry:    taskEntry,
			AssignedUser: assignedUser,
		}

		taskEntries = append(taskEntries, taskEntryReturnWithCreator)
	}

	return taskEntries, nil
}

func deleteTaskEntriesByTeamId(teamId uuid.UUID) error {
	query := database.TaskEntryDeleteByTeamIdQuery
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

func MarkTaskEntryCompleteEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskEntryId := c.Params("taskEntryId")

	taskEntryToMarkComplete, err := getTaskEntryByTaskEntryId(taskEntryId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	task, err := getTaskByTaskId(taskEntryToMarkComplete.TaskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if taskEntryToMarkComplete.Status == "completed" {
		return sendBadRequestResponse(c, err, "Task Entry is already marked complete")
	}

	isUserTheTeamManager, err := isUserTheTeamManager(currentUser.UserId, task.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserTheTeamManager || taskEntryToMarkComplete.AssignedUserId != currentUser.UserId {
		return sendForbiddenResponse(c)
	}

	query := database.TaskEntryMarkCompleteQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(taskEntryId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	// TODO: assign task to next user in the queue

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task Entry Marked Complete",
	})
}

func getTaskEntryByTaskEntryId(taskEntryId string) (*model.TaskEntry, error) {
	query := database.TaskEntryGetByTaskEntryIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	taskEntry := new(model.TaskEntry)

	rows, err := stmt.Query(taskEntryId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return nil, err
		}
	}

	return taskEntry, nil
}

func CancelCurrentTaskEntryEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskEntryId := c.Params("taskEntryId")

	taskEntryToCancel, err := getTaskEntryByTaskEntryId(taskEntryId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	task, err := getTaskByTaskId(taskEntryToCancel.TaskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if taskEntryToCancel.Status == "completed" {
		return sendBadRequestResponse(c, err, "Task Entry is already marked complete")
	}

	isUserTheTeamManager, err := isUserTheTeamManager(currentUser.UserId, task.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserTheTeamManager {
		return sendForbiddenResponse(c)
	}

	query := database.TaskEntryCancelCurrentTaskEntryQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(taskEntryId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	// TODO: assign task to next user in the queue

	return c.SendStatus(fiber.StatusNoContent)
}

func getTaskEntriesByTaskId(taskId uuid.UUID) ([]model.TaskEntryReturnWithAssignedUser, error) {
	query := database.TaskEntriesGetByTaskIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	taskEntries := make([]model.TaskEntryReturnWithAssignedUser, 0)

	rows, err := stmt.Query(taskId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		taskEntry := new(model.TaskEntry)

		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return nil, err
		}

		assignedUser, err := getUserById(taskEntry.AssignedUserId)

		if err != nil {
			return nil, err
		}

		taskEntryReturnWithAssignedUser := model.TaskEntryReturnWithAssignedUser{
			TaskEntry:    taskEntry,
			AssignedUser: assignedUser,
		}

		taskEntries = append(taskEntries, taskEntryReturnWithAssignedUser)
	}

	return taskEntries, nil
}

// func determineNextUserToAssignTaskTo(taskId uuid.UUID) (*model.TaskEntry, error) {
// 	// TODO: implement this function with round robin user assignment

// 	task, err := getTaskByTaskId(taskId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	taskEntries, err := getTaskEntriesByTaskId(taskId)

// 	if err != nil {
// 		return nil, err
// 	}

//   teamMembers, err := getTeamMembersByTeamId(task.TeamId)
// 	return nil, nil
// }
