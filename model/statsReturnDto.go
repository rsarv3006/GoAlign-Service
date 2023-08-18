package model

type StatsReturnDto struct {
	TotalNumberOfTasks           int     `json:"total_number_of_tasks"`
	NumberOfCompletedTaskEntries int     `json:"number_of_completed_task_entries"`
	NumberOfCompletedTasks       int     `json:"number_of_completed_tasks"`
	AverageTasksPerUser          float64 `json:"average_tasks_per_user"`
	TotalNumberOfTaskEntries     int     `json:"total_number_of_task_entries"`
}
