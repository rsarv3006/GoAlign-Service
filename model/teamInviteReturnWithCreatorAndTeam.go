package model

type TeamInviteReturnWithCreator struct {
	*TeamInvite
	Team          TeamReturnWithUsersAndTasks `json:"team"`
	InviteCreator User                        `json:"invite_creator"`
}
