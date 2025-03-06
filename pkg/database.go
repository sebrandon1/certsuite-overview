package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// insertComponentData inserts component details into the dci_components table.
func insertComponentData(db *sql.DB, jobID, commit, createdAt string, totalSuccess, totalFailures, totalErrors, totalSkips int) error {
	if jobID == "" || commit == "" {
		return fmt.Errorf("invalid input: jobID and commit_hash cannot be empty")
	}
	if totalSuccess < 0 || totalFailures < 0 || totalErrors < 0 || totalSkips < 0 {
		return fmt.Errorf("invalid input: totalSuccess=%v, totalFailures=%v, totalErrors=%v, totalSkips=%v", totalSuccess, totalFailures, totalErrors, totalSkips)
	}

	insertQuery := `
        INSERT INTO dci_components (job_id, commit_hash, createdAt, totalSuccess, totalFailures, totalErrors, totalSkips)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE 
        commit_hash = VALUES(commit_hash),
        createdAt = VALUES(createdAt),
        totalSuccess = totalSuccess + VALUES(totalSuccess),
        totalFailures = totalFailures + VALUES(totalFailures),
        totalErrors = totalErrors + VALUES(totalErrors),
        totalSkips = totalSkips + VALUES(totalSkips);
    `
	_, err := db.Exec(insertQuery, jobID, commit, createdAt, totalSuccess, totalFailures, totalErrors, totalSkips)
	return err
}

// insertQuayData inserts a record of Quay image pulls into the aggregated_logs table.
func insertQuayData(db *sql.DB, datetime string, count int, kind string) error {
	log.Printf("Received datetime: %v, count: %v, kind: %v", datetime, count, kind)

	if datetime == "" || kind == "" || count < 0 {
		return fmt.Errorf("invalid input: datetime=%v, kind=%v, count=%d (datetime/kind cannot be empty, count cannot be negative)", datetime, kind, count)

	}

	parsedDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", datetime)
	if err != nil {
		return fmt.Errorf("invalid datetime format: %v, expected YYYY-MM-DD", datetime)
	}
	dateStr := parsedDate.Format("2006-01-02")

	// Define the insert query with ON DUPLICATE KEY UPDATE
	insertQuery := `
    INSERT INTO aggregated_logs (datetime, count, kind)
    VALUES (?, ?, ?)
    ON DUPLICATE KEY UPDATE count = count + VALUES(count);`

	_, err = db.Exec(insertQuery, dateStr, count, kind)
	if err != nil {
		log.Printf("Error executing insert query: %v", err)
	}
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
			datetime DATE NOT NULL,
			count INT NOT NULL,  
			kind VARCHAR(255) NOT NULL,  
			PRIMARY KEY (datetime, kind)		
		);`,

		`CREATE TABLE IF NOT EXISTS dci_components (
			job_id VARCHAR(36) PRIMARY KEY,       
			commit_hash VARCHAR(255) NOT NULL,  
			createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
			totalSuccess INT DEFAULT 0,       
			totalFailures INT DEFAULT 0,      
			totalErrors INT DEFAULT 0,        
			totalSkips INT DEFAULT 0          
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	logrus.Info("Tables created successfully.")
	return nil
}

// chooseDatabase initializes and returns a database connection based on the DB_CHOICE environment variable
func ChooseDatabase() (*sql.DB, error) {
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

func ConnectToAWSDB() (*sql.DB, string, error) {
	logrus.Info("Opening the AWS MySQL database connection...")

	// Fetch the database connection parameters from environment variables.
	DBUsername := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBURL := os.Getenv("DB_URL")
	DBPort := os.Getenv("DB_PORT")

	// Create the initial connection string (without specifying a database).
	initialConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", DBUsername, DBPassword, DBURL, DBPort)
	logrus.Infof("Connecting to MySQL with connection string: %s", initialConnStr)

	// Connect to MySQL without specifying a database.
	db, err := sql.Open("mysql", initialConnStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// Ping the database to ensure the connection is established.
	if err := pingDB(db); err != nil {
		return nil, "", fmt.Errorf("failed to connect to AWS MySQL database: %w", err)
	}

	// Log a successful connection and ping.
	logrus.Info("Successfully connected to AWS MySQL and pinged the database.")

	// Define the database name
	newDBName := "certsuite_usage_db"

	return db, newDBName, nil
}

func initDBAWS() (*sql.DB, error) {
	// Connect to AWS MySQL
	db, newDBName, err := ConnectToAWSDB()
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Create the database if it doesn't exist.
	if err := createDatabase(db, newDBName); err != nil {
		return nil, fmt.Errorf("failed to create database %s: %w", newDBName, err)
	}

	// Close the initial connection and reconnect with the specified database.
	db.Close()
	DBUsername := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBURL := os.Getenv("DB_URL")
	DBPort := os.Getenv("DB_PORT")
	finalConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBURL, DBPort, newDBName)
	logrus.Infof("Reconnecting to MySQL with new database: %s", newDBName)

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

func ConnectToLocalDB() (*sql.DB, error) {
	const (
		rootDSN  = "root:mypassword@tcp(localhost:3306)/"
		dbName   = "certsuite_usage_db"
		finalDSN = "root:mypassword@tcp(localhost:3306)/" + dbName
	)

	logrus.Info("Opening connection to local MySQL server...")

	// Initial connection (to check database existence)
	db, err := sql.Open("mysql", rootDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer db.Close()

	// Check if the database exists
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create database if it does not exist
	if exists == 0 {
		logrus.Infof("Database '%s' does not exist, creating...", dbName)
		if _, err := db.Exec("CREATE DATABASE " + dbName); err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		logrus.Infof("Database '%s' created successfully.", dbName)
	} else {
		logrus.Infof("Database '%s' already exists, skipping creation.", dbName)
	}

	// Connect to the newly ensured database
	db, err = sql.Open("mysql", finalDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database '%s': %w", dbName, err)
	}

	// Verify connection
	if err := pingDB(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MySQL database '%s': %w", dbName, err)
	}

	logrus.Infof("Successfully connected to MySQL database '%s'.", dbName)
	return db, nil
}

func initDBLocal() (*sql.DB, error) {
	// Connect to MySQL
	db, err := ConnectToLocalDB()
	if err != nil {
		return nil, err
	}

	// Create tables in the database
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables in MySQL database: %w", err)
	}

	logrus.Info("Local MySQL database initialized successfully.")
	return db, nil
}
