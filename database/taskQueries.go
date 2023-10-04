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
  creator_id
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
  $11
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
  AND status = 'active'
  AND start_date <= NOW()
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

const TaskUpdateByTaskIdQuery = `
UPDATE tasks 
SET
  task_name = $1,
  notes = $2,
  start_date = $3,
  end_date = $4,
  required_completions_needed = $5,  
  interval_between_windows_count = $6,
  interval_between_windows_unit = $7,
  window_duration_count = $8,
  window_duration_unit = $9
WHERE task_id = $10
RETURNING *
`

const TaskMarkTaskAsCompleteQuery = `
UPDATE tasks
SET status = 'completed'
WHERE task_id = $1
RETURNING *
`

const TaskIncrementCompletionCountQuery = `
UPDATE tasks
SET completion_count = completion_count + 1
WHERE task_id = $1
RETURNING *
`

const TaskGetTasksByTeamIdArrayQuery = `
SELECT * FROM tasks
WHERE team_id = ANY($1);
`
