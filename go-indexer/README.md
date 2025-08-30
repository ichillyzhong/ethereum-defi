# Go Indexer

Go-based backend service for indexing Ethereum DeFi events and providing REST API endpoints.

## Architecture

```
go-indexer/
├── main.go           # Entry point - starts API server and indexer
├── api/              # REST API handlers
│   └── api.go        # TVL endpoint implementation
├── db/               # Database layer
│   └── database.go   # SQLite operations and models
├── indexer/          # Blockchain event indexer
│   └── indexer.go    # Real-time event listener and processor
├── Staking/          # Generated contract bindings
│   └── Staking.go    # Auto-generated from Solidity contract
└── defi_data.db      # SQLite database file
```

## Features

- **Real-time Event Indexing**: Listens to blockchain events via WebSocket
- **TVL Calculation**: Computes Total Value Locked from deposit/withdraw events
- **REST API**: Provides `/api/tvl` endpoint for querying current TVL
- **SQLite Storage**: Persistent storage for all staking events
- **Concurrent Processing**: API server and indexer run simultaneously

## Dependencies

- **Ethereum Client**: `github.com/ethereum/go-ethereum`
- **Web Framework**: `github.com/gin-gonic/gin`
- **Database**: SQLite3 via `github.com/mattn/go-sqlite3`
- **Big Integer**: Go's `math/big` for handling wei amounts

## Usage

### Environment Variables
```bash
export STAKING_CONTRACT_ADDRESS=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
```

### Run the Service
```bash
go run main.go
```

This starts:
1. **API Server** on `localhost:8080`
2. **Event Indexer** connecting to `ws://localhost:8545`

### API Endpoints

#### GET /api/tvl
Returns the current Total Value Locked in wei.

**Response:**
```json
{"tvl":"270000000000000000000"}
```

## Database Schema

### Events Table
```sql
CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user TEXT NOT NULL,
    amount TEXT NOT NULL,
    event_type TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    block_number INTEGER NOT NULL,
    tx_hash TEXT NOT NULL
);
```

## Event Processing

The indexer listens for two event types:

1. **Deposited**: User stakes tokens
   - Increases TVL
   - Event signature: `Deposited(address user, uint256 amount)`

2. **Withdrawn**: User unstakes tokens  
   - Decreases TVL
   - Event signature: `Withdrawn(address user, uint256 amount)`

## TVL Calculation

TVL is calculated by:
```
TVL = SUM(deposit_amounts) - SUM(withdraw_amounts)
```

All amounts are stored and calculated in wei (1 token = 10^18 wei).

## Error Handling

- **Database Connection**: Fails fast if SQLite connection fails
- **Ethereum Connection**: Retries WebSocket connection on disconnect
- **Big Integer Overflow**: Uses `math/big.Int` for safe arithmetic
- **Event Parsing**: Logs errors but continues processing other events

## Development

### Generate Contract Bindings
```bash
abigen --sol ../hardhat/contracts/Staking.sol --pkg Staking --out Staking/Staking.go
```

### Database Queries
```bash
# View all events
sqlite3 defi_data.db "SELECT * FROM events;"

# Count events by type
sqlite3 defi_data.db "SELECT event_type, COUNT(*) FROM events GROUP BY event_type;"

# Calculate TVL manually
sqlite3 defi_data.db "
SELECT 
  SUM(CASE WHEN event_type = 'deposit' THEN CAST(amount AS INTEGER) ELSE 0 END) -
  SUM(CASE WHEN event_type = 'withdraw' THEN CAST(amount AS INTEGER) ELSE 0 END) as tvl
FROM events;"
```

## Configuration

- **Ethereum RPC**: `ws://localhost:8545` (Hardhat node)
- **Database**: `./defi_data.db` (SQLite file)
- **API Port**: `:8080`
- **Contract Address**: Set via `STAKING_CONTRACT_ADDRESS` environment variable
