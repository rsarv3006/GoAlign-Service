package database

const UserTeamMembershipCreateQuery = `
INSERT INTO user_team_memberships (user_id, team_id)
VALUES ($1, $2)
RETURNING id, user_id, team_id, created_at, updated_at
`
