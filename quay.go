package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	quay "github.com/sebrandon1/go-quay/lib"
)

const (
	bearerToken = "wdSLOoT8K3pa4BQhsTnWxLSobTWMNF4X67p4VmU1"
	DateFormat  = "01/02/2006"
)

// fetchQuayData fetches the number of image pulls from Quay.
func fetchQuayData() error {
	quayClient, err := quay.NewClient(bearerToken)
	if err != nil {
		return err
	}

	namespace := "certsuite"
	repository := "redhat-best-practices-for-k8s"
	startDate, endDate := getTodayAndYesterday()

	data, err := quayClient.GetAggregatedLogs(namespace, repository, startDate, endDate)
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
