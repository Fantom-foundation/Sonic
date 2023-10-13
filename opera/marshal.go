package opera

import "encoding/json"

func UpdateRules(src Rules, diff []byte) (res Rules, err error) {
	changed := src.Copy()
	err = json.Unmarshal(diff, &changed)
	if err != nil {
		return src, err
	}
	// protect readonly fields
	res = changed
	res.NetworkID = src.NetworkID
	res.Name = src.Name
	// norma specific override of MinGasPrice by overridden value
	if OverrideMinGasPrice != nil && OverrideMinGasPrice.Sign() > 0 {
		res.Economy.MinGasPrice = OverrideMinGasPrice
	}
	return
}
