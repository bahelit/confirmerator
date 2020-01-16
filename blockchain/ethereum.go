package blockchain

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// getBalance Returns the current balance from the latest block.
func GetBalance(client *ethclient.Client, address common.Address) *big.Float {
	_account := address

	// Passing the block number let's you read the account balance at the time of that block.
	// The block number must be a big.Int.
	//blockNumber := big.NewInt(5532993)
	balance, err := client.BalanceAt(context.Background(), _account, nil)
	if err != nil {
		log.Printf("INFO: failed to get balance for: %v - err: %v", _account, err)
	}
	//fmt.Println(balance)

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	return ethValue
}

func QueryLatestBlockHeader(client *ethclient.Client) (*types.Header, error) {
	// You can call the client's HeaderByNumber to return header information about a block.
	// It'll return the latest block header if you pass nil
	return client.HeaderByNumber(context.Background(), nil)
}

func QueryBlockByNumber(client *ethclient.Client, blockNumber *int64) (*types.Block, error) {
	var tmpBigInt big.Int

	tmpBigInt.SetInt64(*blockNumber)
	block, err := client.BlockByNumber(context.Background(), &tmpBigInt)
	if err != nil {
		time.Sleep(2 * time.Second)
		block, err = client.BlockByNumber(context.Background(), &tmpBigInt)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Failed to get block by header: %v", err)
		}
	}
	log.Println("Block number: ", block.Number())

	return block, nil
}

func QueryBlockByHeader(client *ethclient.Client, header *types.Header) (*types.Block, error) {
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		time.Sleep(2 * time.Second)
		block, err = client.BlockByNumber(context.Background(), header.Number)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Failed to get block by header: %v", err)
		}
	}
	fmt.Println("Block number: ", block.Number())

	return block, nil
}
