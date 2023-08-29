package database

const CreateLoginRequestQuery = `
INSERT INTO login_requests (
  user_id,
  login_request_expiration_date,
  login_request_token
)
VALUES (
  $1, $2, $3
)
RETURNING *;
`

const LoginRequestGetByLoginRequestId = `
SELECT * FROM login_requests
WHERE login_request_id = $1;
`

const LoginRequestMarkAsCompletedQuery = `
UPDATE login_requests
SET login_request_status = 'completed'
WHERE login_request_id = $1;
`
