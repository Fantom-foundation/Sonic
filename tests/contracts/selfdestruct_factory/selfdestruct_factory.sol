// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "selfdestruct/selfdestruct.sol";

contract SelfDestructFactory {
    SelfDestruct lastContract;

    event LogDeployed(address addr);
    event LogContractStorage(uint256 value);

    function create() public payable {
        // construct and transfer value to new contract
        lastContract = (new SelfDestruct){value: msg.value}(
            false, // do not selfdestruct immediately
            false, // beneficiary is not self
            payable(address(0))
        );
        emit LogDeployed(address(lastContract));
    }

    function destroy(address payable beneficiary) public {
        lastContract.destroyContract(false, beneficiary);
    }

    function destroyWithoutBeneficiary() public {
        lastContract.destroyContract(true, payable(address(0)));
    }

    function createAndDestroy(address payable beneficiary) public payable {
        create();
        destroy(beneficiary);
        emit LogContractStorage(lastContract.someData());
    }

    function createAndDestroyWithoutBeneficiary() public payable {
        create();
        destroy(payable(address(lastContract)));
        emit LogContractStorage(lastContract.someData());
    }
}
