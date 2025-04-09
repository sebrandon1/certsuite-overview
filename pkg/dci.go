package pkg

import (
	"fmt"
	"log"
	"strings"

	"github.com/redhat-best-practices-for-k8s/certsuite-overview/config"
	dci "github.com/sebrandon1/go-dci/lib"
)

const (
	daysBackLimit  = 1 // we want display for 1 day
	certsuiteTests = "certsuite-tests_junit.xml"
)

func FetchDciData() error {
	var totalErrors, totalFailures, totalSkips, totalSuccess int

	// Initialize DCI client
	dciClient := dci.NewClient(config.AppConfig.ClientID, config.AppConfig.APISecret)

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

	log.Printf("Fetching DCI data for the last %d days", daysBackLimit)

	// Fetch DCI runs
	runs, err := dciClient.GetJobs(daysBackLimit)
	if err != nil {
		return fmt.Errorf("failed to fetch DCI runs: %w", err)
	}

	log.Printf("Fetched %d DCI runs", len(runs))

	// Store job and component data in the database
	for _, run := range runs {
		for _, job := range run.Jobs {
			// Insert component information into the dci_components table
			for _, component := range job.Components {
				commitHash := "unknown"
				if parts := strings.Split(component.Name, " "); len(parts) > 1 {
					commitHash = parts[1]
				}
				if strings.Contains(component.Name, "cnf-certification-test") || strings.Contains(component.Name, "certsuite") {
					totalErrors = 0
					totalFailures = 0
					totalSkips = 0
					totalSuccess = 0
					for _, result := range job.Results {
						if result.Name == certsuiteTests {
							totalErrors += result.Errors
							totalFailures += result.Failures
							totalSkips += result.Skips
							totalSuccess += result.Success
						}
					}

					log.Println("Inserting DCI component data into the database...")
					log.Printf("Job ID: %s, Commit: %s, CreatedAt: %v, TotalSuccess: %d, TotalFailures: %d, TotalErrors: %d, TotalSkips: %d",
						job.ID, commitHash, job.CreatedAt, totalSuccess, totalFailures, totalErrors, totalSkips)
					log.Println("--------------------")

					if err = insertComponentData(db, job.ID, commitHash, job.CreatedAt, totalSuccess, totalFailures, totalErrors, totalSkips); err != nil {
						log.Printf(
							"Error inserting DCI component entry: Job ID: %s, Commit: %s, CreatedAt: %v, TotalSuccess: %d, TotalFailures: %d, TotalErrors: %d, TotalSkips: %d. Error: %v",
							job.ID, commitHash, job.CreatedAt, totalSuccess, totalFailures, totalErrors, totalSkips, err)
						return fmt.Errorf("failed to insert DCI component data: %w", err)
					}
				}
			}
		}
	}
	log.Println("Successfully fetched and stored DCI data.")
	return nil
}
