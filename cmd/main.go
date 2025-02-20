package main

import (
	"log"

	"github.com/redhat-best-practices-for-k8s/certsuite-overview/config"
	"github.com/spf13/cobra"
	_ "github.com/go-sql-driver/mysql"
)

// Command for 'fetch' action
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch certsuite usage from Quay and DCI",
	Run: func(cmd *cobra.Command, args []string) {
		// Fetch data from Quay and DCI and store it in the database
		if err := FetchCertsuiteUsage(); err != nil {
			log.Fatalf("Failed to fetch certsuite usage: %v", err)
		}
		log.Println("Certsuite usage fetched successfully")
	},
}

// Root command
var rootCmd = &cobra.Command{
	Use:   "certsuite-overview",
	Short: "A CLI to interact with certsuite data",
}

func init() {
	config.LoadConfig()
	rootCmd.AddCommand(fetchCmd)
}

func main() {
	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
