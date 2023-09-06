package model

type TaskEntryReturnWithAssignedUser struct {
	*TaskEntry
	AssignedUser User `json:"assignedUser"`
}
