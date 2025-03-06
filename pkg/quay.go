package pkg

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/redhat-best-practices-for-k8s/certsuite-overview/config"
	quay "github.com/sebrandon1/go-quay/lib"
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

	// Fetch aggregated logs from Quay
	data, err := quayClient.GetAggregatedLogs(config.AppConfig.Namespace, config.AppConfig.Repository, "", "")
	if err != nil {
		return fmt.Errorf("failed to fetch aggregated logs from Quay: %w", err)
	}
	// Loop through the aggregated data and insert it into the database
	for _, aggregated := range data.Aggregated {
		log.Println("Inserting Quay data into the database...")
		log.Printf("Datetime: %s, Count: %d, Kind: %s", aggregated.Datetime, aggregated.Count, aggregated.Kind)
		log.Println("--------------------")
		if err = insertQuayData(db, aggregated.Datetime, aggregated.Count, aggregated.Kind); err != nil {
			log.Printf("Failed to insert Quay data (Datetime: %s, Count: %d, Kind: %s): %v", aggregated.Datetime, aggregated.Count, aggregated.Kind, err)
			return fmt.Errorf("failed to insert Quay data: %w", err)
		}
	}
	log.Println("Successfully fetched and stored Quay data.")
	return nil
}
