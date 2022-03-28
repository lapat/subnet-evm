
const BigNumber = require('bignumber.js');

/* --------------------------- */
const { expect } = require("chai");
const hre = require("hardhat");
//const { ethers, upgrades } = require("hardhat");
var fs = require('fs');
const fetch = require('node-fetch')
async function main(){

  await testEvmFun();
}
async function testEvmFun(){

  var testFactory = await ethers.getContractFactory("TestEvmFunctions");
  console.log("got factory add:"+accounts[0].address)
  var tConnection = await testFactory.connect(accounts[0]);
  console.log("connected")
  var tInstance = await tConnection.deploy();
  console.log("deployed")
  var add = tInstance.address;
  var tx = await tInstance.deployTransaction.wait()
  console.log("add:"+add);
  //var p = tInstance.testCallFun();

}

main()
.then(() => process.exit(0))
.catch(error => {
  console.error(error);
  process.exit(1);
});
