const hre = require("hardhat");

async function main() {
  // Get contract addresses from deployment
  const myTokenAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
  const stakingAddress = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";

  // Get signers
  const [owner, user1, user2] = await hre.ethers.getSigners();

  // Get contract instances
  const MyToken = await hre.ethers.getContractFactory("MyToken");
  const myToken = MyToken.attach(myTokenAddress);

  const Staking = await hre.ethers.getContractFactory("Staking");
  const staking = Staking.attach(stakingAddress);

  console.log("Starting staking tests...");

  // Transfer tokens to users for testing
  const transferAmount = hre.ethers.parseEther("1000");
  console.log("Transferring tokens to users...");
  
  await myToken.transfer(user1.address, transferAmount);
  await myToken.transfer(user2.address, transferAmount);

  // User1 stakes 100 tokens
  const stakeAmount1 = hre.ethers.parseEther("100");
  console.log(`User1 approving ${hre.ethers.formatEther(stakeAmount1)} tokens...`);
  await myToken.connect(user1).approve(stakingAddress, stakeAmount1);
  
  console.log(`User1 staking ${hre.ethers.formatEther(stakeAmount1)} tokens...`);
  await staking.connect(user1).deposit(stakeAmount1);

  // User2 stakes 200 tokens
  const stakeAmount2 = hre.ethers.parseEther("200");
  console.log(`User2 approving ${hre.ethers.formatEther(stakeAmount2)} tokens...`);
  await myToken.connect(user2).approve(stakingAddress, stakeAmount2);
  
  console.log(`User2 staking ${hre.ethers.formatEther(stakeAmount2)} tokens...`);
  await staking.connect(user2).deposit(stakeAmount2);

  // Wait a bit
  await new Promise(resolve => setTimeout(resolve, 2000));

  // User1 withdraws 50 tokens
  const withdrawAmount1 = hre.ethers.parseEther("10");
  console.log(`User1 withdrawing ${hre.ethers.formatEther(withdrawAmount1)} tokens...`);
  await staking.connect(user1).withdraw(withdrawAmount1);

  // User2 withdraws 100 tokens
  const withdrawAmount2 = hre.ethers.parseEther("20");
  console.log(`User2 withdrawing ${hre.ethers.formatEther(withdrawAmount2)} tokens...`);
  await staking.connect(user2).withdraw(withdrawAmount2);

  console.log("Test transactions completed!");
  console.log("Expected final staked amounts:");
  console.log(`User1: ${hre.ethers.formatEther(stakeAmount1 - withdrawAmount1)} ETH`);
  console.log(`User2: ${hre.ethers.formatEther(stakeAmount2 - withdrawAmount2)} ETH`);
  console.log(`Total TVL: ${hre.ethers.formatEther((stakeAmount1 - withdrawAmount1) + (stakeAmount2 - withdrawAmount2))} ETH`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
