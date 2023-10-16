package database

const TaskEntryCreateQuery = `
INSERT INTO task_entries (
  start_date,
  end_date, 
  notes,
  assigned_user_id,
  task_id
)
VALUES (
  $1, $2, $3, $4, $5  
)
RETURNING *
;
`

const TaskEntryDeleteByTeamIdQuery = `
DELETE FROM task_entries
WHERE task_id IN (SELECT task_id FROM tasks WHERE team_id = $1);
`

const TaskEntryGetByTeamIdQuery = `
SELECT * FROM task_entries
WHERE task_id IN (SELECT task_id FROM tasks WHERE team_id = $1);
`

const TaskEntryGetByAssignedUserIdQuery = `
SELECT * FROM task_entries
WHERE assigned_user_id = $1;
`

const TaskEntryGetByTaskEntryIdQuery = `
SELECT * FROM task_entries
WHERE task_entry_id = $1;
`

const TaskEntryMarkCompleteQuery = `
UPDATE task_entries
SET status = 'completed', 
completed_date = CURRENT_TIMESTAMP
WHERE task_entry_id = $1;
`

const TaskEntryCancelCurrentTaskEntryQuery = `
UPDATE task_entries
SET status = 'cancelled'
WHERE task_entry_id = $1;
`

const TaskEntriesGetByTaskIdQuery = `
SELECT * FROM task_entries
WHERE task_id = $1 LIMIT 20;
`

const TaskEntryUpdateAssignedUserIdQuery = `
UPDATE task_entries
SET assigned_user_id = $1
WHERE task_entry_id = $2;
`

const TaskEntryMarkAsCompleteByTaskIdQuery = `
UPDATE task_entries
SET status = 'completed'
WHERE task_id = $1 AND status = 'active';
`

const TaskEntriesGetByTaskIdArrayQuery = `
SELECT * FROM task_entries
WHERE task_id = ANY($1) LIMIT 20;
`

const TaskEntryDeleteByTaskIdQuery = `
DELETE FROM task_entries
WHERE task_id = $1;
`
