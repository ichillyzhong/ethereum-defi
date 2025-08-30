const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

module.exports = buildModule("DeployModule", (m) => {
  // Deploy MyToken first
  const myToken = m.contract("MyToken");

  // Deploy Staking contract with MyToken address as parameter
  const staking = m.contract("Staking", [myToken]);

  return { myToken, staking };
});