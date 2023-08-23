package database

const TaskCreateQuery = `
INSERT INTO tasks (
  task_name,
  notes, 
  start_date,
  end_date,
  required_completions_needed,
  interval_between_windows_count,
  interval_between_windows_unit,
  window_duration_count,
  window_duration_unit,
  team_id,
  creator_id,
  status  
) VALUES (
  $1, 
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12
)
RETURNING *
;`

const TaskDeleteByTeamIdQuery = `
DELETE FROM tasks
WHERE team_id = $1;
`

const TaskGetTasksByTeamIdQuery = `
SELECT * FROM tasks
WHERE team_id = $1;
`

const TaskGetTasksByAssignedUserIdQuery = `
SELECT * FROM tasks
WHERE task_id IN (
  SELECT task_id FROM task_entries
  WHERE assigned_user_id = $1
);
`

const TaskGetTaskByTaskIdQuery = `
SELECT * FROM tasks
WHERE task_id = $1;
`

const TaskDeleteByTaskIdQuery = `
DELETE FROM tasks
WHERE task_id = $1;
`
