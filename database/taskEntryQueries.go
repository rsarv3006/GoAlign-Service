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
);
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
