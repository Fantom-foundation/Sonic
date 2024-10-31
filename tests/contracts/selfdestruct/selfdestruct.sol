// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SelfDestruct {
    uint256 public someData = 123;
    event LogAfterDestruct(uint256 balance);

    constructor(
        bool destroyNow,
        bool selfBeneficiary,
        address payable recipient
    ) payable {
        if (destroyNow) destroyContract(selfBeneficiary, recipient);
    }

    function destroyContract(
        bool selfBeneficiary,
        address payable recipient
    ) public {
        if (selfBeneficiary) recipient = payable(address(this));
        selfdestruct(recipient);

        // This code will not be executed
        emit LogAfterDestruct(address(this).balance);
    }
}
