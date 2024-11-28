package main

import (
	"log"
)

func main() {
	err := fetchCertsuiteUsage()
	if err != nil {
		log.Fatalf("Failed to fetch certsuite usage: %v", err)
	}
}
