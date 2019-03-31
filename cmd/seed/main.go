// Package main seeds the database with mock data for
// development and demonstration purposes.
package main

import (
	"log"

	"github.com/devenney/form3/common"
	"github.com/devenney/form3/utils"
	"github.com/spf13/viper"
)

const (
	mockDataFile = "payments/test_data/mock_data.json"
)

func main() {
	// Initialise our configuration
	common.InitConfig()

	log.Printf("Seeding mock data from [%s] to DynamoDB at %s...", mockDataFile, viper.GetString("db_endpoint"))

	// Seed the mock data
	_, err := utils.InsertMockData(mockDataFile)
	if err != nil {
		log.Fatalf("Failed to insert mock data: %v", err)
	}
}
