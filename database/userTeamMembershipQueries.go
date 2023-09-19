package database

const UserTeamMembershipCreateQuery = `
INSERT INTO user_team_membership (user_id, team_id)
VALUES ($1, $2)
RETURNING user_team_membership_id, user_id, team_id, created_at, updated_at, status
`

const UserTeamMembershipDeleteByTeamIdQuery = `
DELETE FROM user_team_membership
WHERE team_id = $1;
`

const UserTeamMembershipDeleteByUserIdQuery = `
DELETE FROM user_team_membership
WHERE user_id = $1;
`

const UserTeamMembershipGetByTeamIdQuery = `
SELECT * FROM user_team_membership
WHERE team_id = $1;
`

const UserTeamMembershipGetByUserIdAndTeamIdQuery = `
SELECT * FROM user_team_membership
WHERE user_id = $1 AND team_id = $2;
`

const UserTeamMembershipDeleteQueryString = `
DELETE FROM user_team_membership
WHERE user_id = $1 AND team_id = $2;
`

const UserTeamMembershipGetByUserEmailAndTeamIdQueryString = `
select * from user_team_membership 
where team_id = $1 and 
user_id IN (select user_id from users where email = $2);
`
