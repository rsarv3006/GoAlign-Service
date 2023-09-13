package handler

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTask(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskDto := new(model.TaskCreateDto)

	if err := c.BodyParser(taskDto); err != nil {
		return sendBadRequestResponse(c, err, "Error parsing request body")
	}

	taskName := helper.SanitizeInput(taskDto.TaskName)

	if taskName == "" {
		err := errors.New("Task name cannot be empty")
		return sendBadRequestResponse(c, err, "Task name cannot be empty")
	}

	notes := ""
	if taskDto.Notes != nil {
		notes = helper.SanitizeInput(notes)
	}

	if taskDto.StartDate.Before(helper.GetToday()) {
		err := errors.New("Start date cannot be in the past")
		return sendBadRequestResponse(c, err, "Start date cannot be in the past")
	}

	if taskDto.EndDate != nil && taskDto.EndDate.Before(taskDto.StartDate) {
		err := errors.New("End date cannot be before start date")
		return sendBadRequestResponse(c, err, "End date cannot be before start date")
	}

	requiredCompletionsNeeded := taskDto.RequiredCompletionsNeeded
	if requiredCompletionsNeeded != nil && *requiredCompletionsNeeded < 0 {
		err := errors.New("Required completions needed cannot be negative")
		return sendBadRequestResponse(c, err, "Required completions needed cannot be negative")
	} else if requiredCompletionsNeeded == nil {
		// TODO: Is there a better way to do this?
		requiredCompletionsNeeded = new(int)
		*requiredCompletionsNeeded = -1
	}

	if taskDto.IntervalBetweenWindows.IntervalCount < 0 {
		err := errors.New("Interval between windows count cannot be negative")
		return sendBadRequestResponse(c, err, "Interval between windows count cannot be negative")
	}

	if taskDto.WindowDuration.IntervalCount < 0 {
		err := errors.New("Window duration count cannot be negative")
		return sendBadRequestResponse(c, err, "Window duration count cannot be negative")
	}

	if !model.IsValidVariant(string(taskDto.IntervalBetweenWindows.IntervalUnit)) {
		err := errors.New("Interval between windows unit is invalid")
		return sendBadRequestResponse(c, err, "Interval between windows unit is invalid")
	}

	if !model.IsValidVariant(string(taskDto.WindowDuration.IntervalUnit)) {
		err := errors.New("Window duration unit is invalid")
		return sendBadRequestResponse(c, err, "Window duration unit is invalid")
	}

	team, err := getTeamById(taskDto.TeamId)

	if err != nil {
		if err == sql.ErrNoRows {
			err := errors.New("Team does not exist")
			return sendBadRequestResponse(c, err, "Team does not exist")
		} else {
			return sendInternalServerErrorResponse(c, err)
		}
	}

	if team.TeamManagerId != currentUser.UserId {
		err := errors.New("Only the team creator can create tasks for the team")
		return sendBadRequestResponse(c, err, "Only the team creator can create tasks for the team")
	}

	query := database.TaskCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	task := new(model.Task)

	rows, err := stmt.Query(
		taskName,
		notes,
		taskDto.StartDate,
		taskDto.EndDate,
		requiredCompletionsNeeded,
		taskDto.IntervalBetweenWindows.IntervalCount,
		taskDto.IntervalBetweenWindows.IntervalUnit,
		taskDto.WindowDuration.IntervalCount,
		taskDto.WindowDuration.IntervalUnit,
		taskDto.TeamId,
		currentUser.UserId,
	)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	intervalBetweenWindows := model.IntervalObj{}
	windowDuration := model.IntervalObj{}

	for rows.Next() {
		err := rows.Scan(&task.TaskId,
			&task.TaskName,
			&task.Notes,
			&task.StartDate,
			&task.EndDate,
			&task.RequiredCompletionsNeeded,
			&task.CompletionCount,
			&intervalBetweenWindows.IntervalCount,
			&intervalBetweenWindows.IntervalUnit,
			&windowDuration.IntervalCount,
			&windowDuration.IntervalUnit,
			&task.TeamId,
			&task.CreatorId,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Status)
		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}
	}

	task.IntervalBetweenWindows = intervalBetweenWindows
	task.WindowDuration = windowDuration

	endDateInterval := model.IntervalObj{}
	endDateInterval.IntervalCount = task.WindowDuration.IntervalCount
	endDateInterval.IntervalUnit = model.IntervalVariant(task.WindowDuration.IntervalUnit)

	endDate, err := helper.FindDateFromDateAndInterval(task.StartDate, endDateInterval)

	if err != nil {
		return sendBadRequestResponse(c, err, "Error calculating end date")
	}

	if taskDto.EndDate != nil && taskDto.EndDate.Before(endDate) {
		return sendBadRequestResponse(c, err, "End date cannot be before calculated end date")
	}

	taskEntryCreateDto := new(model.TaskEntryCreateDto)
	taskEntryCreateDto.TaskId = task.TaskId
	taskEntryCreateDto.StartDate = task.StartDate
	taskEntryCreateDto.EndDate = endDate
	taskEntryCreateDto.Notes = task.Notes
	taskEntryCreateDto.AssignedUserId = taskDto.AssignedUserId

	taskEntry, err := createTaskEntry(taskEntryCreateDto, currentUser.UserId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	taskObjToReturn := new(model.TaskReturnWithTaskEntries)
	taskObjToReturn.Task = task
	taskObjToReturn.TaskEntries = append(taskObjToReturn.TaskEntries, *taskEntry)
	taskObjToReturn.Creator = *currentUser

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
		"task":    taskObjToReturn,
	})
}

func GetTasksForUserEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	query := database.TaskGetTasksByAssignedUserIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(currentUser.UserId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	tasks := make([]*model.TaskReturnWithTaskEntries, 0)

	for rows.Next() {
		task := new(model.Task)
		intervalBetweenWindows := model.IntervalObj{}
		windowDuration := model.IntervalObj{}

		err := rows.Scan(&task.TaskId,
			&task.TaskName,
			&task.Notes,
			&task.StartDate,
			&task.EndDate,
			&task.RequiredCompletionsNeeded,
			&task.CompletionCount,
			&intervalBetweenWindows.IntervalCount,
			&intervalBetweenWindows.IntervalUnit,
			&windowDuration.IntervalCount,
			&windowDuration.IntervalUnit,
			&task.TeamId,
			&task.CreatorId,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Status)

		task.IntervalBetweenWindows = intervalBetweenWindows
		task.WindowDuration = windowDuration

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		taskEntries, err := getTaskEntriesByTaskId(task.TaskId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		creator, err := getUserById(task.CreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		taskReturnWithTaskEntries := model.TaskReturnWithTaskEntries{
			Task:        task,
			TaskEntries: taskEntries,
			Creator:     creator,
		}

		tasks = append(tasks, &taskReturnWithTaskEntries)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tasks retrieved successfully",
		"tasks":   tasks,
	})
}

func GetTasksByTeamIdEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	teamId, err := uuid.Parse(c.Params("teamId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid teamId")
	}

	isUserInTeam, err := isUserInTeam(currentUser.UserId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserInTeam {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tasks, err := getTasksByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tasks retrieved successfully",
		"tasks":   tasks,
	})
}

func getTasksByTeamId(teamId uuid.UUID) ([]model.TaskReturnWithTaskEntries, error) {
	query := database.TaskGetTasksByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(teamId)

	if err != nil {
		return nil, err
	}

	tasks := []model.TaskReturnWithTaskEntries{}

	for rows.Next() {
		task := new(model.Task)
		intervalBetweenWindows := model.IntervalObj{}
		windowDuration := model.IntervalObj{}

		err := rows.Scan(&task.TaskId,
			&task.TaskName,
			&task.Notes,
			&task.StartDate,
			&task.EndDate,
			&task.RequiredCompletionsNeeded,
			&task.CompletionCount,
			&intervalBetweenWindows.IntervalCount,
			&intervalBetweenWindows.IntervalUnit,
			&windowDuration.IntervalCount,
			&windowDuration.IntervalUnit,
			&task.TeamId,
			&task.CreatorId,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Status)

		if err != nil {
			return nil, err
		}

		task.IntervalBetweenWindows = intervalBetweenWindows
		task.WindowDuration = windowDuration

		taskEntries, err := getTaskEntriesByTaskId(task.TaskId)

		if err != nil {
			return nil, err
		}

		creator, err := getUserById(task.CreatorId)

		if err != nil {
			return nil, err
		}

		taskReturn := model.TaskReturnWithTaskEntries{
			Task:        task,
			Creator:     creator,
			TaskEntries: taskEntries,
		}

		tasks = append(tasks, taskReturn)
	}

	return tasks, nil
}

func deleteTasksByTeamId(teamId uuid.UUID) error {
	query := database.TaskDeleteByTeamIdQuery
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

func isUserTheTeamManager(userId uuid.UUID, teamId uuid.UUID) (bool, error) {
	query := database.TeamGetByTeamIdAndManagerIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	var managerId uuid.UUID

	err = stmt.QueryRow(teamId, userId).Scan(&managerId)

	if err != nil {
		return false, err
	}

	return managerId == userId, nil
}

func DeleteTaskByTaskIdEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskId, err := uuid.Parse(c.Params("taskId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "error parsing taskId")
	}

	task, err := getTaskByTaskId(taskId)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Task not found",
				"error":   err,
			})
		}

		return sendInternalServerErrorResponse(c, err)
	}

	isUserTheTeamManager, err := isUserTheTeamManager(currentUser.UserId, task.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserTheTeamManager {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	err = deleteTaskByTaskId(taskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)

}

func deleteTaskByTaskId(taskId uuid.UUID) error {
	query := database.TaskDeleteByTaskIdQuery
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

func GetTaskEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskId, err := uuid.Parse(c.Params("taskId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "error parsing taskId")
	}

	task, err := getTaskByTaskId(taskId)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Task not found",
			})
		}

		return sendInternalServerErrorResponse(c, err)
	}

	isUserInTeam, err := isUserInTeam(currentUser.UserId, task.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserInTeam {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task retrieved successfully",
		"task":    task,
	})
}

func getTaskByTaskId(taskId uuid.UUID) (*model.Task, error) {
	query := database.TaskGetTaskByTaskIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	task := new(model.Task)

	intervalBetweenWindows := model.IntervalObj{}
	windowDuration := model.IntervalObj{}

	err = stmt.QueryRow(taskId).Scan(&task.TaskId,
		&task.TaskName,
		&task.Notes,
		&task.StartDate,
		&task.EndDate,
		&task.RequiredCompletionsNeeded,
		&task.CompletionCount,
		&intervalBetweenWindows.IntervalCount,
		&intervalBetweenWindows.IntervalUnit,
		&windowDuration.IntervalCount,
		&windowDuration.IntervalUnit,
		&task.TeamId,
		&task.CreatorId,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Status)

	if err != nil {
		return nil, err
	}

	task.IntervalBetweenWindows = intervalBetweenWindows
	task.WindowDuration = windowDuration

	return task, nil
}

func UpdateTaskEndpoint(c *fiber.Ctx) error {
	// TODO: handle updates to assigned user id since that's on the task entry row
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskUpdateDto := new(model.TaskUpdateDto)

	if err := c.BodyParser(taskUpdateDto); err != nil {
		return sendBadRequestResponse(c, err, "error parsing body")
	}

	taskId := taskUpdateDto.TaskId
	taskToUpdate, err := getTaskByTaskId(taskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamId := taskToUpdate.TeamId
	isUserTheTeamManager, err := isUserTheTeamManager(currentUser.UserId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserTheTeamManager {
		return sendForbiddenResponse(c)
	}

	query := database.TaskUpdateByTaskIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(taskUpdateDto.TaskName, taskUpdateDto.Notes, taskUpdateDto.StartDate, taskUpdateDto.EndDate, taskUpdateDto.RequiredCompletionsNeeded, taskUpdateDto.IntervalBetweenWindowsCount, taskUpdateDto.IntervalBetweenWindowsUnit, taskUpdateDto.WindowDurationCount, taskUpdateDto.WindowDurationUnit, taskUpdateDto.TaskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task updated successfully",
	})
}

func incrementTaskCompletionCount(taskId uuid.UUID) error {
	query := database.TaskIncrementCompletionCountQuery
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

func canTaskBeMarkedAsComplete(taskId uuid.UUID) (bool, error) {
	task, err := getTaskByTaskId(taskId)

	if err != nil {
		return false, err
	}

	if task.RequiredCompletionsNeeded == -1 {
		return false, nil
	}

	if task.CompletionCount >= task.RequiredCompletionsNeeded {
		return true, nil
	}

	if task.EndDate != nil && task.EndDate.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

func markTaskAsComplete(taskId uuid.UUID) error {
	query := database.TaskMarkTaskAsCompleteQuery
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
