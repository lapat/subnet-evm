//SPDX-License-Identifier: MIT
pragma solidity >=0.6.2;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./INativeMinter.sol";
import "./INativeTest.sol";
import "hardhat/console.sol";

contract ERC20NativeMinter is ERC20, Ownable {
  // Precompiled Native Minter Contract Address
  address constant MINTER_ADDRESS = 0x0200000000000000000000000000000000000001;
  address constant TEST_FUN_ADDRESS=0x0200000000000000000000000000000000000002;

  // Designated Blackhole Address
  address constant BLACKHOLE_ADDRESS = 0x0100000000000000000000000000000000000000;
  string private constant TOKEN_NAME = "ERC20NativeMinterToken";
  string private constant TOKEN_SYMBOL = "XMPL";

  INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);
//  INativeTest nativeTest = INativeTest(TEST_FUN_ADDRESS);

  event Deposit(address indexed dst, uint256 wad);
  event Mintdrawal(address indexed src, uint256 wad);
  event TestFunction(address indexed src, uint256 wad);
  //event TestReturnNumber();

  constructor(uint256 initSupply) ERC20(TOKEN_NAME, TOKEN_SYMBOL) {
    // Mints INIT_SUPPLY to owner
    _mint(_msgSender(), initSupply);
  }

  // Mints [amount] number of ERC20 token to [to] address.
  function mint(address to, uint256 amount) external onlyOwner {
    _mint(to, amount);
  }

  // Burns [amount] number of ERC20 token from [from] address.
  function burn(address from, uint256 amount) external onlyOwner {
    _burn(from, amount);
  }

  // Swaps [amount] number of ERC20 token for native coin.
  function mintdraw(uint256 wad) external {
    console.log('calling mintdraw - in solidity');

    // Burn ERC20 token first.
    _burn(_msgSender(), wad);
    // Mints [amount] number of native coins (gas coin) to [msg.sender] address.
    // Calls NativeMinter precompile through INativeMinter interface.
    nativeMinter.mintNativeCoin(_msgSender(), wad);
    emit Mintdrawal(_msgSender(), wad);
    console.log('done with mintdraw');

  }

  // Swaps [amount] number of ERC20 token for native coin.
  /*
  function testFunction(uint256 wad) external {
    console.log('calling testFunction - in solidity');
        nativeTest.testFunction(_msgSender(), wad);
    emit TestFunction(_msgSender(), wad);
    console.log('done with testFunction');
  }

  // Swaps [amount] number of ERC20 token for native coin.
  function testReturnNumber() external view returns (uint256){
    console.log('calling testReturnNumber - in solidity');
    // Burn ERC20 token first.
    //_burn(_msgSender(), wad);
    // Mints [amount] number of native coins (gas coin) to [msg.sender] address.
    // Calls NativeMinter precompile through INativeMinter interface.
    //emit TestReturnNumber();
    //console.log('done with testReturnNumber');
    return 12;//nativeTest.testReturnNumber();
  }*/

  // Swaps [amount] number of native gas coins for ERC20 tokens.
  function deposit() external payable {
    // Burn native token by sending to BLACKHOLE_ADDRESS
    payable(BLACKHOLE_ADDRESS).transfer(msg.value);
    // Mint ERC20 token.
    _mint(_msgSender(), msg.value);
    emit Deposit(_msgSender(), msg.value);
  }

  function decimals() public view virtual override returns (uint8) {
    return 18;
  }
}
