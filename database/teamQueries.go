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
