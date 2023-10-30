package handler

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
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

	defer rows.Close()

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

	rows, err := stmt.Query(teamId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userIds := make([]uuid.UUID, 0)
	taskEntries := make([]model.TaskEntry, 0)

	for rows.Next() {
		taskEntry := new(model.TaskEntry)

		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		taskEntries = append(taskEntries, *taskEntry)
		userIds = append(userIds, taskEntry.AssignedUserId)
	}

	users, err := getUsersByIdArray(userIds)

	if err != nil {
		return nil, err
	}

	taskEntriesWithAssignedUsers := make([]model.TaskEntryReturnWithAssignedUser, 0)

	for _, taskEntry := range taskEntries {
		assignedUser := users[taskEntry.AssignedUserId]

		taskEntryWithAssignedUser := model.TaskEntryReturnWithAssignedUser{
			TaskEntry:    &taskEntry,
			AssignedUser: assignedUser,
		}

		taskEntriesWithAssignedUsers = append(taskEntriesWithAssignedUsers, taskEntryWithAssignedUser)
	}

	return taskEntriesWithAssignedUsers, nil
}

func deleteTaskEntriesByTaskId(taskId uuid.UUID) error {
	query := database.TaskEntryDeleteByTaskIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(taskId)

	if err != nil {
		return err
	}

	return nil
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
	currentUser := c.Locals("currentUser").(*model.User)

	taskEntryIdString := c.Params("taskEntryId")
	taskEntryId, err := uuid.Parse(taskEntryIdString)

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid Task Entry Id")
	}

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

	err = incrementTaskCompletionCount(task.TaskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	canTaskBeMarkedAsComplete, err := canTaskBeMarkedAsComplete(task.TaskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if canTaskBeMarkedAsComplete {
		err = markTaskAsComplete(task.TaskId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)

		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Task Marked Complete",
		})
	}

	nextUserId, err, status := determineNextUserToAssignTaskTo(task.TaskId, taskEntryId)

	if err != nil {
		if status == 404 {
			return sendNotFoundResponse(c, err.Error())
		} else if status == 400 {
			return sendBadRequestResponse(c, err, err.Error())
		} else {
			return sendInternalServerErrorResponse(c, err)
		}
	}

	newTaskEntry, err := createTaskEntryFromPreviousTaskEntry(taskEntryId, *nextUserId, currentUser.UserId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Task Entry Marked Complete",
		"newTaskEntry": newTaskEntry,
	})

}

func getTaskEntryByTaskEntryId(taskEntryId uuid.UUID) (*model.TaskEntry, error) {
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

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return nil, err
		}
	}

	return taskEntry, nil
}

func CancelCurrentTaskEntryEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	taskEntryIdString := c.Params("taskEntryId")

	taskEntryId, err := uuid.Parse(taskEntryIdString)

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid Task Entry Id")
	}

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

	_, err = stmt.Exec(taskEntryId)

	defer stmt.Close()

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	canTaskBeMarkedAsComplete, err := canTaskBeMarkedAsComplete(task.TaskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if canTaskBeMarkedAsComplete {
		err := markTaskAsComplete(task.TaskId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}
	} else {
		userIdForNextTaskEntry, err, status := determineNextUserToAssignTaskTo(task.TaskId, taskEntryId)

		if err != nil {
			if status == 404 {
				return sendNotFoundResponse(c, err.Error())
			} else if status == 400 {
				return sendBadRequestResponse(c, err, err.Error())
			} else {
				return sendInternalServerErrorResponse(c, err)
			}
		}

		_, err = createTaskEntryFromPreviousTaskEntry(taskEntryId, *userIdForNextTaskEntry, currentUser.UserId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func getTaskEntriesByTaskIdArray(taskIds []uuid.UUID) (map[uuid.UUID][]model.TaskEntryReturnWithAssignedUser, error) {
	rows, err := database.POOL.Query(
		context.Background(),
		database.TaskEntriesGetByTaskIdArrayQuery,
		pq.Array(taskIds))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	taskEntries := make([]model.TaskEntry, 0)

	userIds := make([]uuid.UUID, 0)

	for rows.Next() {
		taskEntry := new(model.TaskEntry)

		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return nil, err
		}

		taskEntries = append(taskEntries, *taskEntry)
		userIds = append(userIds, taskEntry.AssignedUserId)
	}

	users, err := getUsersByIdArray(userIds)

	if err != nil {
		return nil, err
	}

	taskEntriesMap := make(map[uuid.UUID][]model.TaskEntryReturnWithAssignedUser)

	for _, taskEntry := range taskEntries {
		newTaskEntry := taskEntry
		taskEntryWithAssignedUser := model.TaskEntryReturnWithAssignedUser{
			TaskEntry:    &newTaskEntry,
			AssignedUser: users[taskEntry.AssignedUserId],
		}

		taskEntriesMap[taskEntry.TaskId] = append(taskEntriesMap[taskEntry.TaskId], taskEntryWithAssignedUser)
	}

	return taskEntriesMap, nil
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

	defer rows.Close()

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

func determineNextUserToAssignTaskTo(taskId uuid.UUID, taskEntryId uuid.UUID) (*uuid.UUID, error, int) {
	task, err := getTaskByTaskId(taskId)

	if err != nil {
		return nil, err, 500
	}

	teamMembers, err := getUsersByTeamId(task.TeamId)

	if err != nil {
		return nil, err, 500
	}

	if len(teamMembers) == 1 {
		return &teamMembers[0].UserId, nil, 201
	}

	taskEntries, err := getTaskEntriesByTaskId(taskId)

	currentTaskEntry := model.TaskEntryReturnWithAssignedUser{}

	for _, taskEntry := range taskEntries {
		if taskEntry.TaskEntryId == taskEntryId {
			currentTaskEntry = taskEntry
		}
	}

	if currentTaskEntry.Status == "" {
		return nil, errors.New("Task Entry not found"), 404
	}

	if err != nil {
		return nil, err, 500
	}

	memberTaskEntryCountMap := make(map[uuid.UUID]int)

	for _, teamMember := range teamMembers {
		memberTaskEntryCountMap[teamMember.UserId] = 0
	}

	for _, taskEntry := range taskEntries {
		if taskEntry.Status == "completed" {
			memberTaskEntryCountMap[taskEntry.AssignedUser.UserId]++
		}
	}

	delete(memberTaskEntryCountMap, currentTaskEntry.AssignedUserId)

	var minTaskEntryCount = 0
	var minTaskEntryCountUserId uuid.UUID

	for userId, taskEntryCount := range memberTaskEntryCountMap {
		if taskEntryCount <= minTaskEntryCount {
			minTaskEntryCount = taskEntryCount
			minTaskEntryCountUserId = userId
		}
	}

	return &minTaskEntryCountUserId, nil, 201
}

func createTaskEntryFromPreviousTaskEntry(
	previousTaskEntryId uuid.UUID,
	nextAssignedUserId uuid.UUID,
	currentUserId uuid.UUID) (*model.TaskEntryReturnWithAssignedUser, error) {
	previousTaskEntry, err := getTaskEntryByTaskEntryId(previousTaskEntryId)

	if err != nil {
		return nil, err
	}

	task, err := getTaskByTaskId(previousTaskEntry.TaskId)

	if err != nil {
		return nil, err
	}

	startDate, err := helper.FindDateFromDateAndInterval(previousTaskEntry.EndDate, task.IntervalBetweenWindows)

	if err != nil {
		return nil, err
	}

	endDate, err := helper.FindDateFromDateAndInterval(startDate, task.WindowDuration)

	if err != nil {
		return nil, err
	}

	taskEntryCreateDto := model.TaskEntryCreateDto{
		StartDate:      startDate,
		EndDate:        endDate,
		Notes:          "",
		AssignedUserId: nextAssignedUserId,
		TaskId:         previousTaskEntry.TaskId,
	}

	newTaskEntry, err := createTaskEntry(&taskEntryCreateDto, currentUserId)

	return newTaskEntry, err
}

func updateTaskEntryAssignedUserId(taskEntryId uuid.UUID, userId uuid.UUID) error {
	query := database.TaskEntryUpdateAssignedUserIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(userId, taskEntryId)

	if err != nil {
		return err
	}

	return nil
}

func markTaskEntryAsCompleteByTaskId(taskId uuid.UUID) error {
	query := database.TaskEntryMarkAsCompleteByTaskIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(taskId)

	if err != nil {
		return err
	}

	return nil
}
