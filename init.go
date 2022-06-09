package wallet

import "math/big"

var ChainID *big.Int

func init() {
	ChainID = big.NewInt(0x539)
}
