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
