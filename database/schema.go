package database

import "log"

func CreateUserTable() {
	log.Println("Creating users table")
	_, err := DB.Query(`CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY default gen_random_uuid(),
  user_name VARCHAR(50) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  is_email_verified BOOLEAN NOT NULL DEFAULT FALSE
);
`)

	if err != nil {
		panic(err)
	}
}

func CreateTeamTable() {
	log.Println("Creating teams table")
	_, err := DB.Query(`CREATE TABLE IF NOT EXISTS teams (
  team_id UUID NOT NULL PRIMARY KEY default gen_random_uuid(),
  team_name VARCHAR(255) NOT NULL,
  creator_user_id UUID NOT NULL,
  status VARCHAR(255) NOT NULL DEFAULT 'active',
  team_manager_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);`)

	if err != nil {
		panic(err)
	}
}

func CreateUserTeamMembershipTable() {
	log.Println("Creating user_team_membership table")
	_, err := DB.Query(`CREATE TABLE IF NOT EXISTS user_team_membership (
  user_team_membership_id UUID NOT NULL PRIMARY KEY default gen_random_uuid(),
  user_id UUID NOT NULL,
  team_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  status VARCHAR(255) NOT NULL DEFAULT 'active'
);`)

	if err != nil {
		panic(err)
	}
}

func CreateTeamSettingsTable() {
	log.Println("Creating team_settings table")
	_, err := DB.Query(`
    CREATE TABLE IF NOT EXISTS team_settings (
      team_settings_id UUID NOT NULL PRIMARY KEY default gen_random_uuid(),
      team_id UUID NOT NULL,
      user_id UUID NOT NULL,
      status VARCHAR(255) NOT NULL DEFAULT 'active',
      can_all_members_add_tasks BOOLEAN NOT NULL DEFAULT FALSE
    );
  `)

	if err != nil {
		panic(err)
	}
}
