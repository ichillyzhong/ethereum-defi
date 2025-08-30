const {
  loadFixture,
} = require("@nomicfoundation/hardhat-toolbox/network-helpers");
const { expect } = require("chai");

describe("MyToken and Staking", function () {
  // Fixture to deploy both contracts
  async function deployTokenAndStakingFixture() {
    const [owner, user1, user2] = await ethers.getSigners();

    // Deploy MyToken
    const MyToken = await ethers.getContractFactory("MyToken");
    const myToken = await MyToken.deploy();

    // Deploy Staking contract with MyToken address
    const Staking = await ethers.getContractFactory("Staking");
    const staking = await Staking.deploy(myToken.target);

    // Transfer some tokens to users for testing
    const transferAmount = ethers.parseEther("1000");
    await myToken.transfer(user1.address, transferAmount);
    await myToken.transfer(user2.address, transferAmount);

    return { myToken, staking, owner, user1, user2, transferAmount };
  }

  describe("MyToken Deployment", function () {
    it("Should set the right token name and symbol", async function () {
      const { myToken } = await loadFixture(deployTokenAndStakingFixture);

      expect(await myToken.name()).to.equal("MyToken");
      expect(await myToken.symbol()).to.equal("MTK");
    });

    it("Should assign the total supply to the owner", async function () {
      const { myToken, owner } = await loadFixture(deployTokenAndStakingFixture);
      const expectedSupply = ethers.parseEther("10000");

      expect(await myToken.totalSupply()).to.equal(expectedSupply);
      expect(await myToken.balanceOf(owner.address)).to.be.greaterThan(0);
    });

    it("Should have 18 decimals", async function () {
      const { myToken } = await loadFixture(deployTokenAndStakingFixture);

      expect(await myToken.decimals()).to.equal(18);
    });
  });

  describe("Staking Contract Deployment", function () {
    it("Should set the right token address", async function () {
      const { myToken, staking } = await loadFixture(deployTokenAndStakingFixture);

      expect(await staking.token()).to.equal(myToken.target);
    });

    it("Should set the right owner", async function () {
      const { staking, owner } = await loadFixture(deployTokenAndStakingFixture);

      expect(await staking.owner()).to.equal(owner.address);
    });
  });

  describe("Staking Functionality", function () {
    describe("Deposit", function () {
      it("Should allow users to deposit tokens", async function () {
        const { myToken, staking, user1 } = await loadFixture(deployTokenAndStakingFixture);
        const depositAmount = ethers.parseEther("100");

        // Approve staking contract to spend tokens
        await myToken.connect(user1).approve(staking.target, depositAmount);

        // Deposit tokens
        await expect(staking.connect(user1).deposit(depositAmount))
          .to.emit(staking, "Deposited")
          .withArgs(user1.address, depositAmount, await ethers.provider.getBlock("latest").then(b => b.timestamp + 1));

        // Check staked balance
        expect(await staking.stakedBalances(user1.address)).to.equal(depositAmount);
      });

      it("Should revert when depositing zero amount", async function () {
        const { staking, user1 } = await loadFixture(deployTokenAndStakingFixture);

        await expect(staking.connect(user1).deposit(0))
          .to.be.revertedWith("Amount must be greater than zero");
      });

      it("Should revert when insufficient allowance", async function () {
        const { staking, user1 } = await loadFixture(deployTokenAndStakingFixture);
        const depositAmount = ethers.parseEther("100");

        // Don't approve, should fail
        await expect(staking.connect(user1).deposit(depositAmount))
          .to.be.reverted;
      });

      it("Should update user balance correctly after multiple deposits", async function () {
        const { myToken, staking, user1 } = await loadFixture(deployTokenAndStakingFixture);
        const firstDeposit = ethers.parseEther("100");
        const secondDeposit = ethers.parseEther("50");

        // First deposit
        await myToken.connect(user1).approve(staking.target, firstDeposit);
        await staking.connect(user1).deposit(firstDeposit);

        // Second deposit
        await myToken.connect(user1).approve(staking.target, secondDeposit);
        await staking.connect(user1).deposit(secondDeposit);

        const expectedTotal = firstDeposit + secondDeposit;
        expect(await staking.stakedBalances(user1.address)).to.equal(expectedTotal);
      });
    });

    describe("Withdraw", function () {
      it("Should allow users to withdraw staked tokens", async function () {
        const { myToken, staking, user1 } = await loadFixture(deployTokenAndStakingFixture);
        const depositAmount = ethers.parseEther("100");
        const withdrawAmount = ethers.parseEther("50");

        // First deposit tokens
        await myToken.connect(user1).approve(staking.target, depositAmount);
        await staking.connect(user1).deposit(depositAmount);

        const initialBalance = await myToken.balanceOf(user1.address);

        // Withdraw tokens
        await expect(staking.connect(user1).withdraw(withdrawAmount))
          .to.emit(staking, "Withdrawn")
          .withArgs(user1.address, withdrawAmount, await ethers.provider.getBlock("latest").then(b => b.timestamp + 1));

        // Check balances
        expect(await staking.stakedBalances(user1.address)).to.equal(depositAmount - withdrawAmount);
        expect(await myToken.balanceOf(user1.address)).to.equal(initialBalance + withdrawAmount);
      });

      it("Should revert when withdrawing zero amount", async function () {
        const { staking, user1 } = await loadFixture(deployTokenAndStakingFixture);

        await expect(staking.connect(user1).withdraw(0))
          .to.be.revertedWith("Amount must be greater than zero");
      });

      it("Should revert when withdrawing more than staked balance", async function () {
        const { myToken, staking, user1 } = await loadFixture(deployTokenAndStakingFixture);
        const depositAmount = ethers.parseEther("100");
        const withdrawAmount = ethers.parseEther("200");

        // Deposit some tokens first
        await myToken.connect(user1).approve(staking.target, depositAmount);
        await staking.connect(user1).deposit(depositAmount);

        // Try to withdraw more than deposited
        await expect(staking.connect(user1).withdraw(withdrawAmount))
          .to.be.revertedWith("Insufficient staked balance");
      });

      it("Should revert when user has no staked balance", async function () {
        const { staking, user2 } = await loadFixture(deployTokenAndStakingFixture);
        const withdrawAmount = ethers.parseEther("100");

        await expect(staking.connect(user2).withdraw(withdrawAmount))
          .to.be.revertedWith("Insufficient staked balance");
      });
    });

    describe("Multiple Users", function () {
      it("Should handle multiple users staking independently", async function () {
        const { myToken, staking, user1, user2 } = await loadFixture(deployTokenAndStakingFixture);
        const user1Deposit = ethers.parseEther("100");
        const user2Deposit = ethers.parseEther("200");

        // User1 deposits
        await myToken.connect(user1).approve(staking.target, user1Deposit);
        await staking.connect(user1).deposit(user1Deposit);

        // User2 deposits
        await myToken.connect(user2).approve(staking.target, user2Deposit);
        await staking.connect(user2).deposit(user2Deposit);

        // Check individual balances
        expect(await staking.stakedBalances(user1.address)).to.equal(user1Deposit);
        expect(await staking.stakedBalances(user2.address)).to.equal(user2Deposit);
      });

      it("Should allow partial withdrawals for multiple users", async function () {
        const { myToken, staking, user1, user2 } = await loadFixture(deployTokenAndStakingFixture);
        const depositAmount = ethers.parseEther("100");
        const withdrawAmount = ethers.parseEther("30");

        // Both users deposit
        await myToken.connect(user1).approve(staking.target, depositAmount);
        await staking.connect(user1).deposit(depositAmount);

        await myToken.connect(user2).approve(staking.target, depositAmount);
        await staking.connect(user2).deposit(depositAmount);

        // Both users withdraw partially
        await staking.connect(user1).withdraw(withdrawAmount);
        await staking.connect(user2).withdraw(withdrawAmount);

        const expectedRemaining = depositAmount - withdrawAmount;
        expect(await staking.stakedBalances(user1.address)).to.equal(expectedRemaining);
        expect(await staking.stakedBalances(user2.address)).to.equal(expectedRemaining);
      });
    });
  });
});
