package database

func CreateProductTable() {
	DB.Query(`CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    amount integer,
    name text UNIQUE,
    description text,
    category text NOT NULL
)
`)
}

func CreateUserTable() {
	DB.Query(`CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY,
  user_name VARCHAR(50), 
  email VARCHAR(100),
  is_active BOOLEAN
  is_email_verified BOOLEAN
);
`)
}

func CreateTeamTable() {
	DB.Query(`CREATE TABLE IF NOT EXISTS teams (
  team_id UUID,
  team_name VARCHAR(255),
  creator_user_id UUID,
  status VARCHAR(255),
  team_manager_id UUID,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (team_id)
);`)
}

func CreateUserTeamMembershipTable() {
	DB.Query(`CREATE TABLE IF NOT EXISTS user_team_membership (
  user_id UUID,
  team_id UUID,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  status VARCHAR(255),
  PRIMARY KEY (user_id, team_id)
);`)
}
