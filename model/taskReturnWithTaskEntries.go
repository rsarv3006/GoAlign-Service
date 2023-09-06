package model

type TaskReturnWithTaskEntries struct {
	*Task
	Creator     User                              `json:"creator"`
	TaskEntries []TaskEntryReturnWithAssignedUser `json:"taskEntries"`
}
