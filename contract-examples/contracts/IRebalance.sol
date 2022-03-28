//SPDX-License-Identifier: MIT
pragma solidity >=0.6.2;

interface IRebalance{
  function getCoingeckoPrice(address addr) external view returns (uint256);
}
