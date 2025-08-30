# Ethereum DeFi Protocol
A complete DeFi staking protocol with real-time event indexing and TVL analytics

## Project Structure
- `hardhat/` - Smart contracts (MyToken ERC20 & Staking contract)
- `go-indexer/` - Go backend service for indexing blockchain events and API

## Setup
npm install --save-dev hardhat@^2.22.15

## Testing the Complete DeFi System

### Prerequisites
Make sure you have the following installed:
- Node.js and npm
- Go (1.19+)
- SQLite3 (for database)

### Step 1: Start Hardhat Node
Open Terminal 1:
```bash
cd hardhat
npx hardhat node
```
This starts a local Ethereum network on `localhost:8545` with pre-funded test accounts.

### Step 2: Deploy Smart Contracts
Open Terminal 2:
```bash
cd hardhat
npx hardhat ignition deploy ./ignition/modules/deploy.js --network localhost
```
This deploys MyToken and Staking contracts. Note the deployed addresses:
- MyToken: `0x5FbDB2315678afecb367f032d93F642f64180aa3`
- Staking: `0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512`

### Step 3: Start API Server & Event Indexer
Open Terminal 3:
```bash
cd go-indexer
STAKING_CONTRACT_ADDRESS=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512 go run main.go
```
This starts both the REST API server on `localhost:8080` and the blockchain event indexer concurrently. The indexer captures staking events in real-time while the API serves TVL data.

### Step 4: Generate Test Staking Events
Open Terminal 4:
```bash
cd hardhat
npx hardhat run scripts/test-staking.js --network localhost
```
This script:
- Transfers tokens to test users
- Executes deposit transactions (User1: 100 tokens, User2: 200 tokens)
- Executes withdraw transactions (User1: 10 tokens, User2: 20 tokens)
- Expected final TVL: 270 tokens

### Step 5: Verify Results
**Important**: Make sure all services from Steps 1-4 are running before verification.

Check the database directly:
```bash
sqlite3 go-indexer/defi_data.db "SELECT * FROM events;"
```

Check the API endpoint:
```bash
curl http://localhost:8080/api/tvl
```

### Expected Output
- **Indexer logs**: "Successfully inserted deposit/withdraw event for user..."
- **Database**: Should contain 4 events (2 deposits, 2 withdrawals)
- **API response**: `{"tvl":"270000000000000000000"}` (270 tokens in wei)

### Troubleshooting
If API returns `{"error":"Failed to get TVL"}`:
1. Ensure all services are running (Steps 1-4)
2. Check if events were captured: `sqlite3 go-indexer/defi_data.db "SELECT COUNT(*) FROM events;"`
3. If no events, restart the indexer and run the test script again
