// SPDX-License-Identifier: MIT
pragma solidity ^0.8.8;

contract BurnGas {
    uint256[200] public data;

    function burnGas() public {
        for (uint256 i = 0; i < data.length; i++) {
            data[i] = i;
        }
    }
}
