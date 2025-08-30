package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Event 是我们数据库表的 Go 结构体表示
type Event struct {
	User        string
	Amount      *big.Int
	EventType   string
	Timestamp   time.Time
	BlockNumber uint64
	TxHash      string
}

// DB 结构体用于管理数据库连接
type DB struct {
	*sql.DB
}

// NewDB 创建并返回一个新的数据库连接实例
func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database!")
	return &DB{db}, nil
}

// CreateTables 创建存储事件的表
func (d *DB) CreateTables() error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_address TEXT NOT NULL,
		amount TEXT NOT NULL,
		event_type TEXT NOT NULL,
		block_number INTEGER NOT NULL,
		tx_hash TEXT NOT NULL UNIQUE,
		timestamp DATETIME NOT NULL
	);
	`
	_, err := d.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create events table: %v", err)
	}
	log.Println("Events table created or already exists.")
	return nil
}

// InsertEvent 将一个事件数据插入到数据库
func (d *DB) InsertEvent(event *Event) error {
	query := `
	INSERT OR IGNORE INTO events (user_address, amount, event_type, block_number, tx_hash, timestamp)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := d.Exec(query, event.User, event.Amount.String(), event.EventType, event.BlockNumber, event.TxHash, event.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert event: %v", err)
	}
	return nil
}

// GetTotalValueLocked 查询总锁仓价值 (TVL)
func (d *DB) GetTotalValueLocked() (*big.Int, error) {
	// Get all deposit events
	depositQuery := `SELECT amount FROM events WHERE event_type = 'deposit'`
	depositRows, err := d.Query(depositQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query deposits: %v", err)
	}
	defer depositRows.Close()

	totalDeposits := big.NewInt(0)
	for depositRows.Next() {
		var amountStr string
		if err := depositRows.Scan(&amountStr); err != nil {
			return nil, fmt.Errorf("failed to scan deposit amount: %v", err)
		}
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			continue // Skip invalid amounts
		}
		totalDeposits.Add(totalDeposits, amount)
	}

	// Get all withdraw events
	withdrawQuery := `SELECT amount FROM events WHERE event_type = 'withdraw'`
	withdrawRows, err := d.Query(withdrawQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query withdrawals: %v", err)
	}
	defer withdrawRows.Close()

	totalWithdrawals := big.NewInt(0)
	for withdrawRows.Next() {
		var amountStr string
		if err := withdrawRows.Scan(&amountStr); err != nil {
			return nil, fmt.Errorf("failed to scan withdrawal amount: %v", err)
		}
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			continue // Skip invalid amounts
		}
		totalWithdrawals.Add(totalWithdrawals, amount)
	}

	// Calculate TVL = deposits - withdrawals
	tvl := new(big.Int).Sub(totalDeposits, totalWithdrawals)
	return tvl, nil
}
