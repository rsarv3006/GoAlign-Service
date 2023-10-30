package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/donutsahoy/yourturn-fiber/config"
)

var DB *sql.DB
var POOL *pgxpool.Pool

func Connect() error {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		fmt.Println("Error parsing str to int")
	}

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"), config.Config("SSL_MODE"))

	POOL, err = pgxpool.New(context.Background(), connectionString)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	// defer POOL.Close()

	DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if !fiber.IsChild() {
		fmt.Println("Creating tables if needed...")
		CreateUserTable()
		CreateTeamTable()
		CreateUserTeamMembershipTable()
		CreateTeamSettingsTable()
		CreateTaskTable()
		CreateTaskEntryTable()
		CreateTeamInviteTable()
		CreateAppLogsTable()
		CreateLoginRequestsTable()
	}

	fmt.Println("Connection Opened to Database")
	return nil
}
