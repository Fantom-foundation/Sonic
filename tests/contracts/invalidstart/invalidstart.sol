// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract InvalidStart {
    function createWithInvalidCode() public {
        assembly {
            mstore(0, calldataload(0))
            if iszero(create(0, 0, calldatasize())) {
                revert(0, 0)
            }
        }
    }

    function create2WithInvalidCode() public {
        assembly {
            mstore(0, calldataload(0))
            if iszero(create2(0, 0, calldatasize(), 0)) {
                revert(0, 0)
            }
        }
    }

    // This function creates a contract with no code and attempts to transfer to it.
    function createEmptyContractAndTransferToIt() public {
        bytes memory bytecode = hex""; // empty code
        address newContract;
        bool result;

        assembly {
            newContract := create(0, add(bytecode, 0x20), mload(bytecode))
            result := call(gas(), newContract, 1, 0, 0, 0, 0)
            if iszero(result) {
                revert(0, 0)
            }
        }
    }
}
