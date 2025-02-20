package pkg

import (
	"fmt"
	"log"
	"strings"

	dci "github.com/sebrandon1/go-dci/lib"
	"github.com/redhat-best-practices-for-k8s/certsuite-overview/config"
	
)

const (
	daysBackLimit = 1 // we want display for 1 day
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

	// Fetch DCI runs
	runs, err := dciClient.GetJobs(daysBackLimit)
	if err != nil {
		return fmt.Errorf("failed to fetch DCI runs: %w", err)
	}

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
					for _, result := range job.Results {
						totalErrors += result.Errors
						totalFailures += result.Failures
						totalSkips += result.Skips
						totalSuccess += result.Success
					}
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
