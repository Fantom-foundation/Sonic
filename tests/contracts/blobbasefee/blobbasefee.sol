// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BlobBaseFee {
    event CurrentBlobBaseFee(uint256 fee);
    event LogBlobGasUsed(uint256 gasUsed);

    function logCurrentBlobBaseFee() public {
        emit CurrentBlobBaseFee(block.blobbasefee);
    }

    function getBlobBaseFee() public view returns (uint256) {
        return block.blobbasefee;
    }
}
