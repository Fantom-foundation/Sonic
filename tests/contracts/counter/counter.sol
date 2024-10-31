// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Counter {
    int private count = 0;

    function incrementCounter() public {
        count += 1;
    }

    function getCount() public view returns (int) {
        return count;
    }
}

contract MiniSfc {

    // Parameters
    uint256 public targetGasPowerPerSecond;
    uint256 public counterweight;

    // Result
    uint256 public minGasPrice;

    function _sealEpoch_minGasPrice(uint256 epochDuration, uint256 epochGas) internal {
        // change minGasPrice proportionally to the difference between target and received epochGas
        uint256 targetEpochGas = epochDuration * targetGasPowerPerSecond + 1;
        uint256 gasPriceDeltaRatio = epochGas * GP.SCALING / targetEpochGas;
        // scale down the change speed (estimate gasPriceDeltaRatio ^ (epochDuration / counterweight))
        gasPriceDeltaRatio = (epochDuration * gasPriceDeltaRatio + counterweight * GP.SCALING) / (epochDuration + counterweight);
        // limit the max/min possible delta in one epoch
        gasPriceDeltaRatio = GP.trimGasPriceChangeRatio(gasPriceDeltaRatio);

        // apply the ratio
        uint256 newMinGasPrice = minGasPrice * gasPriceDeltaRatio / GP.SCALING;
        // limit the max/min possible minGasPrice
        newMinGasPrice = GP.trimMinGasPrice(newMinGasPrice);
        // apply new minGasPrice
        minGasPrice = newMinGasPrice;
    }


    function _sealEpoch_minGasPrice_rewrite(uint256 epochDuration, uint256 epochGas) internal {
        // change minGasPrice proportionally to the difference between target and received epochGas
        // scale down the change speed (estimate gasPriceDeltaRatio ^ (epochDuration / counterweight))
        uint256 gasPriceDeltaRatio = 
            (epochDuration * ((epochGas * GP.SCALING) / (epochDuration * targetGasPowerPerSecond + 1)) + counterweight * GP.SCALING) 
            / (epochDuration + counterweight);
  
  
        // limit the max/min possible delta in one epoch
        gasPriceDeltaRatio = GP.trimGasPriceChangeRatio(gasPriceDeltaRatio);

        // apply the ratio
        uint256 newMinGasPrice = minGasPrice * gasPriceDeltaRatio / GP.SCALING;
        // limit the max/min possible minGasPrice
        newMinGasPrice = GP.trimMinGasPrice(newMinGasPrice);
        // apply new minGasPrice
        minGasPrice = newMinGasPrice;
    }
}



library GP {

    uint256 constant SCALING = 1e18;

    function trimGasPriceChangeRatio(uint256 x) internal pure returns (uint256) {
        if (x > SCALING * 105 / 100) {
            return SCALING * 105 / 100;
        }
        if (x < SCALING * 95 / 100) {
            return SCALING * 95 / 100;
        }
        return x;
    }

    function trimMinGasPrice(uint256 x) internal pure returns (uint256) {
        if (x > 1000000 * 1e9) {
            return 1000000 * 1e9;
        }
        if (x < 1e9) {
            return 1e9;
        }
        return x;
    }

    function initialMinGasPrice() internal pure returns (uint256) {
        return 100 * 1e9;
    }
}