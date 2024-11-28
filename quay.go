package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shimritproj/certsuite-overview/config"

	quay "github.com/sebrandon1/go-quay/lib"
)

const (
	DateFormat = "01/02/2006"
)

// fetchQuayData fetches the number of image pulls from Quay.
func fetchQuayData() error {
	cfg, err := config.LoadConfig("config/config.json") // Ensure correct path to config.json
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	quayClient, err := quay.NewClient(cfg.BearerToken)
	if err != nil {
		return err
	}

	startDate, endDate := getTodayAndYesterday()

	data, err := quayClient.GetAggregatedLogs(cfg.Namespace, cfg.Repository, startDate, endDate)
	if err != nil {
		return err
	}

	// Initialize the database connection
	db, err := initDB()
	if err != nil {
		log.Printf("Error initializing database")
		return err
	}
	defer db.Close()

	// Loop through the aggregated data and insert it into the database
	for _, aggregated := range data.Aggregated {
		err = insertQuayData(db, aggregated.Datetime, aggregated.Count, aggregated.Kind)
		if err != nil {
			log.Printf("Failed to insert Quay data: %v", err)
			return err
		}
	}
	return nil
}

// getTodayAndYesterday returns today's date and yesterday's date (24 hours before) in string format
func getTodayAndYesterday() (string, string) {
	// Get the current time
	today := time.Now()

	// Format today's date as a string
	todayStr := today.Format(DateFormat)

	// Calculate yesterday's date by subtracting 24 hours
	yesterday := today.Add(-24 * time.Hour)
	yesterdayStr := yesterday.Format(DateFormat)

	return todayStr, yesterdayStr
}
