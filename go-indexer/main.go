package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ichillyzhong/ethereum-defi/go-indexer/api"
	"github.com/ichillyzhong/ethereum-defi/go-indexer/db"
	"github.com/ichillyzhong/ethereum-defi/go-indexer/indexer"
)

var (
	DbConnStr         = "./defi_data.db"
	ethereumClientURL = "ws://localhost:8545"
	stakingContract   = os.Getenv("STAKING_CONTRACT_ADDRESS")
)

func main() {

	// Start indexer
	go indexer.Run(ethereumClientURL, stakingContract)

	// Connect to database
	dbClient, err := db.NewDB(DbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbClient.Close()

	// Create database tables
	if err := dbClient.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Setup API routes
	router := api.SetupRouter(dbClient)

	// Run API server
	log.Println("Starting API server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
