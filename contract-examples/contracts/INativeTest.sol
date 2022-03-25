//SPDX-License-Identifier: MIT
pragma solidity >=0.6.2;
import "./IAllowList.sol";

interface INativeTest is IAllowList {
  // Mint [amount] number of native coins and send to [addr]
  //function mintNativeCoin(address addr, uint256 amount) external;
  function testFunction(address addr, uint256 amount) external;
}
