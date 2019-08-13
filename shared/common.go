// Generic functions that can be used across projects
package shared

import (
	"regexp"

	"github.com/ethereum/go-ethereum/common"
)

type Symbol string

const (
	SymbBitcoin  = Symbol("BTC")
	SymbEthereum = Symbol("ETH")
)

func (symbol *Symbol) ToString() string {
	return string(*symbol)
}

// IsValidAddress validate hex address
func IsValidAddress(address interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := address.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}
