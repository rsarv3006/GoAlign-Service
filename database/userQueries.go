package database

const UserDeleteUserQuery = `
DELETE FROM users
WHERE user_id = $1;
`

const UserGetUserByIdQuery = `
SELECT * FROM users
WHERE user_id = $1;
`

const UserGetByIdArrayQuery = `
SELECT * FROM users
WHERE user_id = ANY($1);
`

const UserGetUserByEmailQuery = `
SELECT * FROM users
WHERE email = $1;
`

const UserGetUsersByTeamIdQuery = `
SELECT * FROM users
WHERE user_id IN (
  SELECT user_id FROM user_team_membership
  WHERE team_id = $1
);
`

const UserCreateUserQuery = `
INSERT INTO users (username, email) VALUES ($1, $2) RETURNING *`

const UserGetUsersByIdArrayQuery = `
SELECT * FROM users
WHERE user_id = ANY($1);
`

const UserGetUsersByTeamIdArrayQuery = `
SELECT u.*, m.team_id
FROM users u
JOIN user_team_membership m ON u.user_id = m.user_id
WHERE m.team_id = ANY($1);
`

const UserDeleteUserLoginRequestsByUserIdQuery = `
DELETE FROM login_requests
WHERE user_id = $1;
`
