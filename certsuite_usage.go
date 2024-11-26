package main

import "fmt"

// fetchCertsuiteUsage integrates data from Quay and DCI.
func fetchCertsuiteUsage() error {
	if err := fetchQuayData(); err != nil {
		return fmt.Errorf("error fetching Quay data: %w", err)
	}
	if err := fetchDciData(); err != nil {
		return fmt.Errorf("error fetching DCI data: %w", err)
	}
	return nil
}
