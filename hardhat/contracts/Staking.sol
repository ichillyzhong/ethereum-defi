// SPDX-License-Identifier: MIT
  pragma solidity ^0.8.20;

  import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
  import "@openzeppelin/contracts/access/Ownable.sol";

  contract Staking is Ownable {
      IERC20 public token;

      event Deposited(address indexed user, uint256 amount, uint256 timestamp);
      event Withdrawn(address indexed user, uint256 amount, uint256 timestamp);

      mapping(address => uint256) public stakedBalances;

      constructor(address _token) Ownable(msg.sender) {
          token = IERC20(_token);
      }

      function deposit(uint256 amount) external {
          require(amount > 0, "Amount must be greater than zero");
          token.transferFrom(msg.sender, address(this), amount);
          stakedBalances[msg.sender] += amount;
          emit Deposited(msg.sender, amount, block.timestamp);
      }

      function withdraw(uint256 amount) external {
          require(amount > 0, "Amount must be greater than zero");
          require(stakedBalances[msg.sender] >= amount, "Insufficient staked balance");
          stakedBalances[msg.sender] -= amount;
          token.transfer(msg.sender, amount);
          emit Withdrawn(msg.sender, amount, block.timestamp);
      }
  }