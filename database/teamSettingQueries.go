package database

const TeamSettingsCreateQuery = `
INSERT INTO team_settings (
  team_id,
  can_all_members_add_tasks
)
VALUES (
  $1, $2
)
RETURNING team_settings_id, team_id, can_all_members_add_tasks;
`

const TeamSettingsDeleteByTeamIdQuery = `
DELETE FROM team_settings
WHERE team_id = $1;
`

const TeamSettingsUpdateQuery = `
UPDATE team_settings
SET can_all_members_add_tasks = $1
WHERE team_id = $2
RETURNING team_settings_id, team_id, can_all_members_add_tasks;
`

const TeamSettingsGetByTeamIdQuery = `
SELECT * FROM team_settings
WHERE team_id = $1;
`
