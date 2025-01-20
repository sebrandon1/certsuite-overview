package pkg

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // for SQLite
	"github.com/sirupsen/logrus"
)

// initDBLocal initializes a local SQLite database and ensures required tables exist.
func initDBLocal() (*sql.DB, error) {
	logrus.Info("Opening the local SQLite database connection...")

	db, err := sql.Open("sqlite3", "certsuite_usage.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	if err := pingDB(db); err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables in SQLite database: %w", err)
	}

	logrus.Info("Local SQLite database initialized successfully.")
	return db, nil
}

// initDBAWS initializes a MySQL database connection using AWS credentials.
func initDBAWS() (*sql.DB, error) {
	logrus.Info("Opening the AWS MySQL database connection...")

	// Fetch the database connection parameters from environment variables.
	DBUsername := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBURL := os.Getenv("DB_URL")
	DBPort := os.Getenv("DB_PORT")

	// Create the connection string for initial connection.
	initialConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", DBUsername, DBPassword, DBURL, DBPort)
	logrus.Infof("Connecting to MySQL with connection string: %s", initialConnStr)

	// Connect to MySQL without specifying a database.
	db, err := sql.Open("mysql", initialConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// Ping the database to ensure the connection is established.
	if err := pingDB(db); err != nil {
		return nil, fmt.Errorf("failed to connect to AWS MySQL database: %w", err)
	}

	// Log a successful connection and ping.
	logrus.Info("Successfully connected to AWS MySQL and pinged the database.")
	
	// Create the new database if it doesn't exist.
	newDBName := "certsuite_usage_db"
	if err := createDatabase(db, newDBName); err != nil {
		return nil, fmt.Errorf("failed to create database %s: %w", newDBName, err)
	}

	// Update the connection string to include the new database.
	finalConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBURL, DBPort, newDBName)
	logrus.Infof("Reconnecting to MySQL with new database: %s", newDBName)
	
	// Close the initial connection and establish a new connection with the database specified.
	db.Close()
	db, err = sql.Open("mysql", finalConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to reconnect to MySQL with database %s: %w", newDBName, err)
	}

	// Ping the database again to ensure the connection is established.
	if err := pingDB(db); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", newDBName, err)
	}

	// Log a successful reconnection and ping with the new database.
	logrus.Infof("Successfully connected to database '%s' and pinged the database.", newDBName)

	// Create tables in the new database.
	logrus.Info("Creating tables in the database...")
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables in database %s: %w", newDBName, err)
	}

	logrus.Info("AWS MySQL database initialized successfully.")
	return db, nil
}
// insertComponentData inserts component details into the dci_components table.
func insertComponentData(db *sql.DB, job_id, commit, createdAt string, totalSuccess, totalFailures, totalErrors, totalSkips int) error {
	insertQuery := `
		INSERT OR REPLACE INTO dci_components (job_id, commit_hash, createdAt, totalSuccess, totalFailures, totalErrors, totalSkips)
		VALUES (?, ?, ?, ?, ?, ?, ?);
    `
	_, err := db.Exec(insertQuery, job_id, commit, createdAt, totalSuccess, totalFailures, totalErrors, totalSkips)
	return err
}

// insertQuayData inserts a record of Quay image pulls into the database.
func insertQuayData(db *sql.DB, datetime string, count int, kind string) error {
	insertQuery := `
        INSERT OR REPLACE INTO aggregated_logs (datetime, count, kind) 
        VALUES (?, ?, ?);
    `
	_, err := db.Exec(insertQuery, datetime, count, kind)
	return err
}

// pingDB verifies the database connection.
func pingDB(db *sql.DB) error {
	logrus.Info("Pinging the database to verify connection...")
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("database ping failed: %w", err)
	}
	logrus.Info("Database connection verified successfully.")
	return nil
}

// createDatabase creates a new database if it doesn't already exist.
func createDatabase(db *sql.DB, dbName string) error {
	logrus.Infof("Checking if database %s exists...", dbName)
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database %s: %w", dbName, err)
	}
	logrus.Infof("Database %s created or already exists.", dbName)
	return nil
}

// createTables creates the required tables if they do not exist.
func createTables(db *sql.DB) error {
	logrus.Info("Creating tables if they do not exist...")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS aggregated_logs (
			id CHAR(36) PRIMARY KEY,  -- Store UUID as CHAR(36)
			datetime TEXT,
			count INTEGER,
			kind TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS dci_components (
			job_id CHAR(36) PRIMARY KEY,       -- Store UUID as CHAR(36)
			commit_hash TEXT NOT NULL,         -- Commit hash
			createdAt TIMESTAMP NOT NULL,      -- Job creation timestamp
			totalSuccess INTEGER DEFAULT 0,   -- Number of successful results
			totalFailures INTEGER DEFAULT 0,  -- Number of failed results
			totalErrors INTEGER DEFAULT 0,    -- Number of errors
			totalSkips INTEGER DEFAULT 0      -- Number of skipped results
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			db.Close()
			return fmt.Errorf("failed to execute table creation query: %w", err)
		}
	}

	logrus.Info("All required tables created successfully.")
	return nil
}

// chooseDatabase initializes and returns a database connection based on the DB_CHOICE environment variable
func chooseDatabase() (*sql.DB, error) {
	dbChoice := os.Getenv("DB_CHOICE") // Expecting "local" or "aws"
	var db *sql.DB
	var err error

	if dbChoice == "aws" {
		db, err = initDBAWS()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize AWS database: %w", err)
		}
	} else {
		db, err = initDBLocal()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize local database: %w", err)
		}
	}

	return db, nil
}
