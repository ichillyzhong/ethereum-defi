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
	// 连接到以太坊客户端
	client, err := ethclient.Dial(ethereumClientURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// 连接到数据库
	dbConnStr := "./defi_data.db"
	dbClient, err := db.NewDB(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbClient.Close()
	dbClient.CreateTables()

	// 实例化合约
	contractAddress := common.HexToAddress(stakingContract)
	contract, err := Staking.NewStaking(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	// 创建事件过滤器
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

			// 尝试解析 Deposit 事件
			depositedEvent, err := contract.ParseDeposited(vLog)
			if err == nil {
				processEvent(dbClient, client, vLog, "deposit", depositedEvent.User, depositedEvent.Amount)
				continue
			}

			// 尝试解析 Withdraw 事件
			withdrawnEvent, err := contract.ParseWithdrawn(vLog)
			if err == nil {
				processEvent(dbClient, client, vLog, "withdraw", withdrawnEvent.User, withdrawnEvent.Amount)
				continue
			}
		}
	}
}

// processEvent 处理并存储事件到数据库
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
