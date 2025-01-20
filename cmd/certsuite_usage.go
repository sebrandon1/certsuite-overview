package main

import (
	"fmt"

	"github.com/redhat-best-practices-for-k8s/certsuite-overview/pkg"
)

// FetchCertsuiteUsage integrates data from Quay and DCI.
func FetchCertsuiteUsage() error {
	if err := pkg.FetchQuayData(); err != nil {
		return fmt.Errorf("error fetching Quay data: %w", err)
	}
	if err := pkg.FetchDciData(); err != nil {
		return fmt.Errorf("error fetching DCI data: %w", err)
	}
	return nil
}
