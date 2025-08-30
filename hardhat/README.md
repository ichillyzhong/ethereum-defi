# Ethereum DeFi Staking Project

This project demonstrates a DeFi staking system built with Hardhat. It includes an ERC20 token contract and a staking contract that allows users to stake tokens and earn rewards.

## Contracts

### MyToken.sol
An ERC20 token contract that serves as the base token for the staking system:
- **Name**: MyToken (MTK)
- **Initial Supply**: 10,000 tokens (minted to deployer)
- **Decimals**: 18
- Built using OpenZeppelin's ERC20 implementation

### Staking.sol
A staking contract that allows users to deposit and withdraw tokens:
- **Deposit**: Users can stake their MyToken tokens
- **Withdraw**: Users can withdraw their staked tokens at any time
- **Balance Tracking**: Maintains individual staked balances for each user
- **Events**: Emits Deposited and Withdrawn events for tracking
- **Access Control**: Uses OpenZeppelin's Ownable for admin functions

## Usage

Try running some of the following tasks:

```shell
npx hardhat help
npx hardhat test
REPORT_GAS=true npx hardhat test
npx hardhat node
npx hardhat ignition deploy ./ignition/modules/deploy.js --network local
```

## Generate Go Bindings
```shell
cp artifacts/contracts/Staking.sol/Staking.json Staking.json
jq '.abi' Staking.json > StakingABI.json
abigen --abi StakingABI.json --pkg Staking --out Staking.go
mv Staking.go go-indexer/Staking/
```
