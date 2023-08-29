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
