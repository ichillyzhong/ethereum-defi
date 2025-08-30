package indexer

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ichillyzhong/ethereum-defi/go-indexer/Staking"
	"github.com/ichillyzhong/ethereum-defi/go-indexer/db"
)

func Run(ethereumClientURL, stakingContract string) {
	// Connect to Ethereum client
	client, err := ethclient.Dial(ethereumClientURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Connect to database
	dbConnStr := "./defi_data.db"
	dbClient, err := db.NewDB(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbClient.Close()
	dbClient.CreateTables()

	// Instantiate contract
	contractAddress := common.HexToAddress(stakingContract)
	contract, err := Staking.NewStaking(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	// Create event filter
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	fmt.Println("Listening for events...")
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case vLog := <-logs:
			fmt.Println("Received a new log!")

			// Try to parse Deposit event
			depositedEvent, err := contract.ParseDeposited(vLog)
			if err == nil {
				processEvent(dbClient, client, vLog, "deposit", depositedEvent.User, depositedEvent.Amount)
				continue
			}

			// Try to parse Withdraw event
			withdrawnEvent, err := contract.ParseWithdrawn(vLog)
			if err == nil {
				processEvent(dbClient, client, vLog, "withdraw", withdrawnEvent.User, withdrawnEvent.Amount)
				continue
			}
		}
	}
}

// processEvent processes and stores events to database
func processEvent(dbClient *db.DB, client *ethclient.Client, vLog types.Log, eventType string, user common.Address, amount *big.Int) {
	header, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(vLog.BlockNumber)))
	if err != nil {
		log.Printf("Failed to get block header for block %d: %v", vLog.BlockNumber, err)
		return
	}

	event := &db.Event{
		User:        user.Hex(),
		Amount:      amount,
		EventType:   eventType,
		Timestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber: vLog.BlockNumber,
		TxHash:      vLog.TxHash.Hex(),
	}

	if err := dbClient.InsertEvent(event); err != nil {
		log.Printf("Failed to insert %s event: %v", eventType, err)
	} else {
		log.Printf("Successfully inserted %s event for user %s", eventType, user.Hex())
	}
}
