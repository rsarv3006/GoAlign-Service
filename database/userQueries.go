package database

const UserDeleteUserQuery = `
DELETE FROM users
WHERE user_id = $1;
`

const UserGetUserByIdQuery = `
SELECT * FROM users
WHERE user_id = $1;
`
