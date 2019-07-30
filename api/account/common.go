package account

import "fmt"

type Chain int16

const (
	ChainBitcoin = iota
	ChainEthereum
	ChainEthereumClassic
	ChainBitcoinCash
	ChainCallisto
	ChainRavenCoin
)

func StringToChain(symbol string) Chain {
	fmt.Println(symbol)
	switch symbol {
	case "BTC":
		return ChainBitcoin
	case "ETH":
		return ChainEthereum
	case "ETC":
		return ChainEthereumClassic
	case "BCH":
		return ChainBitcoinCash
	case "CLO":
		return ChainCallisto
	case "RVN":
		return ChainRavenCoin
	default:
		return ChainEthereum
	}
}
