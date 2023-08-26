package model

type TeamSettingsUpdateDto struct {
	CanAllMembersAddTasks *bool `json:"can_all_members_add_tasks,omitempty"`
}
