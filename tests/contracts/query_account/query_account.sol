// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract QueryAccount {
    function getBalance(address addr) public view returns (uint256) {
        return addr.balance;
    }

    function getCodeSize(address addr) public view returns (uint256) {
        uint256 codeSize;
        assembly {
            codeSize := extcodesize(addr)
        }
        return codeSize;
    }
}
