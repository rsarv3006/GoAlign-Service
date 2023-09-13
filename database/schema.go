package database

import "log"

func CreateUserTable() {
	log.Println("Creating users table")
	_, err := DB.Query(`CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY default gen_random_uuid(),
  username VARCHAR(50) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
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
  creator_user_id UUID NOT NULL REFERENCES users(user_id),
  status VARCHAR(255) NOT NULL DEFAULT 'active',
  team_manager_id UUID NOT NULL REFERENCES users(user_id),
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
  user_id UUID NOT NULL REFERENCES users(user_id),
  team_id UUID NOT NULL REFERENCES teams(team_id),
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
      team_id UUID NOT NULL REFERENCES teams(team_id),
      can_all_members_add_tasks BOOLEAN NOT NULL DEFAULT FALSE
    );
  `)

	if err != nil {
		panic(err)
	}
}

func CreateTaskTable() {
	log.Println("Creating tasks table")
	_, err := DB.Query(`
CREATE TABLE IF NOT EXISTS tasks (
  task_id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  task_name varchar NOT NULL,
  notes text NOT NULL DEFAULT '',
  start_date timestamptz NOT NULL,
  end_date timestamptz NULL DEFAULT NULL,
  required_completions_needed integer NOT NULL DEFAULT -1, 
  completion_count integer NOT NULL DEFAULT 0,
  interval_between_windows_count integer NOT NULL,
  interval_between_windows_unit varchar NOT NULL,
  window_duration_count integer NOT NULL,
  window_duration_unit varchar NOT NULL,
  team_id uuid NOT NULL REFERENCES teams(team_id),
  creator_id uuid NOT NULL REFERENCES users(user_id),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  status varchar NOT NULL DEFAULT 'active'
);
  `)

	if err != nil {
		panic(err)
	}
}

func CreateTaskEntryTable() {
	log.Println("Creating task_entries table")
	_, err := DB.Query(`
  CREATE TABLE IF NOT EXISTS task_entries (
  task_entry_id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  start_date timestamptz NOT NULL,
  end_date timestamptz NOT NULL,
  notes text NOT NULL,
  assigned_user_id uuid NOT NULL REFERENCES users(user_id),
  status varchar(50) NOT NULL DEFAULT 'active',
  completed_date timestamptz,
  task_id uuid NOT NULL REFERENCES tasks(task_id)
);`)

	if err != nil {
		panic(err)
	}
}

func CreateTeamInviteTable() {
	log.Println("Creating team_invites table")
	_, err := DB.Query(`
  CREATE TABLE IF NOT EXISTS team_invites (
  team_invite_id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  team_id uuid NOT NULL REFERENCES teams(team_id),
  email varchar(100) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  status varchar(50) NOT NULL DEFAULT 'pending',
  invite_creator_id uuid NOT NULL REFERENCES users(user_id)
);`)

	if err != nil {
		panic(err)
	}
}

func CreateAppLogsTable() {
	log.Println("Creating app_logs table")
	_, err := DB.Query(`
  CREATE TABLE IF NOT EXISTS app_logs (
  app_log_id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  log_message text NOT NULL,
  log_level varchar(50) NOT NULL,
  log_date timestamptz NOT NULL DEFAULT NOW(),
  log_data jsonb,
  user_id uuid 
);`)

	if err != nil {
		panic(err)
	}
}

func CreateLoginRequestsTable() {
	log.Println("Creating login_requests table")
	_, err := DB.Query(`
  CREATE TABLE IF NOT EXISTS login_requests (
  login_request_id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(user_id),
  login_request_date timestamptz NOT NULL DEFAULT NOW(),
  login_request_expiration_date timestamptz NOT NULL,
  login_request_token varchar(255) NOT NULL,
  login_request_status varchar(50) NOT NULL DEFAULT 'pending'
);`)

	if err != nil {
		panic(err)
	}
}
