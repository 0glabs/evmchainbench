// SPDX-License-Identifier: cc0-1.0 
pragma solidity ^0.8.0;

contract Incrementer {
    uint256 public count = 0;

    function increment() public {
        count += 1;
    }
}
