package model

type TeamReturnWithUsersAndTasks struct {
	*Team
	Users []User                      `json:"users"`
	Tasks []TaskReturnWithTaskEntries `json:"tasks"`
}
