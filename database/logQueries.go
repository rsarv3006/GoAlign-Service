package database

const LogCreateQueryWithJsonAndUserId = `
INSERT INTO app_logs (
  log_message,
  log_level,
  log_data,
  user_id
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *
;`
