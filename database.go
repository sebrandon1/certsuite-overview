package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // for SQLite
)

// initDB initializes the SQLite database and creates the table for storing aggregated logs.
func initDB() (*sql.DB, error) {
	log.Println("Opening the database connection...")

	db, err := sql.Open("sqlite3", "certsuite_usage.db")
	if err != nil {
		log.Printf("Failed to open database: %v\n", err)
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	log.Println("Pinging the database to verify connection...")
	if err = db.Ping(); err != nil {
		db.Close()
		log.Printf("Failed to ping database: %v\n", err)
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Creating tables if they do not exist...")

	// Create the aggregated_logs table if it doesn't exist
	createQuayTableQuery := `
        CREATE TABLE IF NOT EXISTS aggregated_logs (
            id TEXT PRIMARY KEY,  -- Store UUID as TEXT (string)
            datetime TEXT,
            count INTEGER,
            kind TEXT
        );
    `
	_, err = db.Exec(createQuayTableQuery)
	if err != nil {
		db.Close()
		log.Printf("Failed to create aggregated_logs table: %v\n", err)
		return nil, fmt.Errorf("failed to create aggregated_logs table: %v", err)
	}

	// Create the dci_jobs table if it doesn't exist
	createDciTableQuery := `
    	CREATE TABLE IF NOT EXISTS dci_components (
        	job_id TEXT PRIMARY KEY,  -- Store UUID as TEXT()
            commit_hash TEXT NOT NULL,     -- Commit hash
            createdAt DATETIME NOT NULL, -- Job creation timestamp
            totalSuccess INTEGER DEFAULT 0,  -- Number of successful results
            totalFailures INTEGER DEFAULT 0, -- Number of failed results
            totalErrors INTEGER DEFAULT 0,   -- Number of errors
            totalSkips INTEGER DEFAULT 0     -- Number of skipped results
    	);
	`
	_, err = db.Exec(createDciTableQuery)
	if err != nil {
		db.Close()
		log.Printf("Failed to create dci_jobs table: %v\n", err)
		return nil, fmt.Errorf("failed to create dci_jobs table: %v", err)
	}

	log.Println("Database initialized successfully.")
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
