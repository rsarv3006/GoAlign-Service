package database

const TeamCreateQueryString = `
  INSERT INTO teams (
    team_name,
    creator_user_id,
    team_manager_id
  )  
  VALUES ($1, $2, $3)
  RETURNING *;
`

const TeamGetByUserIdQueryString = `
  SELECT * FROM teams WHERE team_id IN (
    SELECT team_id FROM user_team_membership WHERE user_id = $1
  );
`

const TeamDeleteQueryString = `
DELETE FROM teams
WHERE team_id = $1;
`

const TeamGetByIdQueryString = `
  SELECT * FROM teams WHERE team_id = $1;
`

const TeamGetByTeamManagerIdQueryString = `
  SELECT * FROM teams WHERE team_manager_id = $1;
`

const TeamGetByTeamIdAndManagerIdQuery = `
 select team_manager_id from teams where team_id = $1 and team_manager_id = $2;
`

const TeamUpdateTeamManagerQueryString = `
  UPDATE teams
  SET team_manager_id = $1
  WHERE team_id = $2
  RETURNING *;
`

const TeamGetByIdsQueryString = `
select * from teams where team_id = ANY($1);
`
