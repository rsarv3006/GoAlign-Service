package database

const UserTeamMembershipCreateQuery = `
INSERT INTO user_team_membership (user_id, team_id)
VALUES ($1, $2)
RETURNING user_team_membership_id, user_id, team_id, created_at, updated_at, status
`
