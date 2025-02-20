package pkg

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	quay "github.com/sebrandon1/go-quay/lib"
	"github.com/redhat-best-practices-for-k8s/certsuite-overview/config"
)

const (
	DateFormat = "01/02/2006"
)

// fetchQuayData fetches the number of image pulls from Quay.
func FetchQuayData() error {
	// Initialize database connection
	db, err := ChooseDatabase()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Failed to close database connection: %v", closeErr)
		}
	}()

	// Initialize Quay client
	quayClient, err := quay.NewClient(config.AppConfig.BearerToken)
	if err != nil {
		return fmt.Errorf("failed to initialize Quay client: %w", err)
	}

	// Get date range
	startDate, endDate := getTodayAndYesterday()

	// Fetch aggregated logs from Quay
	data, err := quayClient.GetAggregatedLogs(config.AppConfig.Namespace, config.AppConfig.Repository, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to fetch aggregated logs from Quay: %w", err)
	}

	// Loop through the aggregated data and insert it into the database
	for _, aggregated := range data.Aggregated {
		if err = insertQuayData(db, aggregated.Datetime, aggregated.Count, aggregated.Kind); err != nil {
			log.Printf("Failed to insert Quay data (Datetime: %s, Count: %d, Kind: %s): %v", aggregated.Datetime, aggregated.Count, aggregated.Kind, err)
			return fmt.Errorf("failed to insert Quay data: %w", err)
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

	return yesterdayStr,todayStr 
}
