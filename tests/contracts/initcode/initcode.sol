// SPDX-License-Identifier: MIT
pragma solidity ^0.8.8;

contract InitCode {
    event LogCost(uint256 cost);
    
    function createContractWith(uint codeSize) public {
        uint256 before = gasleft();
        assembly ("memory-safe") {
            pop(create(0, 0, codeSize)) 
        }
        emit LogCost(before - gasleft());
    }
}