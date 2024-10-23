// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Coinbase {
    event LogCoinbase(address fee);

    function logCoinbase() public payable {
        emit LogCoinbase(block.coinbase);
    }

    function getCoinbase() public view returns (address) {
        return block.coinbase;
    }
}
