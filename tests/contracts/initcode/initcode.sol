// SPDX-License-Identifier: MIT
pragma solidity ^0.8.8;

contract InitCode {
    event LogCost(uint256 cost);
    
    function createContractWith(uint codeSize) public {
        // save gas before attempting to create
        uint256 before = gasleft();
        assembly {
            // top level expresssions are not supposed to return a value, but create does.
            // so we need to call pop to catch this value.
            pop(create(0, 0, codeSize))
        }
        // report how much gas was used during create call
        emit LogCost(before - gasleft());
    }

    function create2ContractWith(uint codeSize) public {
        // save gas before attempting to create
        uint256 before = gasleft();
        uint256 result = 0;
        assembly  {
            // top level expresssions are not supposed to return a value, but create does.
            // so we need to assign its return value to catch it.
            result := create2(0, 0, codeSize, 0)
        }
        // report how much gas was used during create call
        emit LogCost(before - gasleft());
    }

    // this function is only used to measure the cost of calling a function but not creating a contract
    function measureGasAndAssign(uint codeSize) public {
        uint256 before = gasleft();
        uint256 result = 0;
        assembly  {
            result := codeSize
        }
        emit LogCost(before - gasleft());
    }
}