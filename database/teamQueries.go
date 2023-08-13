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
