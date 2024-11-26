package main

import (
	"log"
	"strings"

	dci "github.com/sebrandon1/go-dci/lib"
)

const (
	daysBackLimit = 1 // we want display for 1 day
)

func fetchDciData() error {
	var totalErrors, totalFailures, totalSkips, totalSuccess int

	// Initialize DCI client
	clientID := "remoteci/8f6c5d9a-3ca6-4f37-bb59-1aa15d0fa43e"
	apiSecret := "z9F9GeEAZvD8rtqucsmHarfnGMHRt4dzsr063syZv5wxxHRzhfEEoDX6MZ6e5yB0"
	dciClient := dci.NewClient(clientID, apiSecret)

	// Initialize the database
	db, err := initDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Fetch DCI runs
	runs, err := dciClient.GetJobs(daysBackLimit)
	if err != nil {
		return err
	}

	// Store job and component data in the database
	for _, run := range runs {
		for _, job := range run.Jobs {
			// Insert component information into the dci_components table
			for _, component := range job.Components {
				commit_hash := "unknown"
				if parts := strings.Split(component.Name, " "); len(parts) > 1 {
					commit_hash = parts[1]
				}
				if strings.Contains(component.Name, "cnf-certification-test") || strings.Contains(component.Name, "certsuite") {
					for _, result := range job.Results {
						totalErrors += result.Errors
						totalFailures += result.Failures
						totalSkips += result.Skips
						totalSuccess += result.Success
					}
					err = insertComponentData(db, job.ID, commit_hash, job.CreatedAt, totalSuccess, totalFailures, totalErrors, totalSkips)
					if err != nil {
						log.Printf(
							"Error inserting DCI component entry: Job ID: %s, Commit: %s, CreatedAt: %v, TotalSuccess: %d, TotalFailures: %d, TotalErrors: %d, TotalSkips: %d. Error: %v",
							job.ID, commit_hash, job.CreatedAt, totalSuccess, totalFailures, totalErrors, totalSkips, err,
						)
						return err
					}
				}
			}
		}

	}

	return nil
}
