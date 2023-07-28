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
	if res.Economy.OverrideMinGasPrice != nil {
		res.Economy.MinGasPrice = res.Economy.OverrideMinGasPrice
	}
	return
}
