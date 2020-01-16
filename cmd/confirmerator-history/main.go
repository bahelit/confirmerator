package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"

	"github.com/bahelit/confirmerator/blockchain"
	ethereum_history "github.com/bahelit/confirmerator/cmd/confirmerator-history/ethereum-history"
	"github.com/bahelit/confirmerator/database/postgres"
)

var (
	pgClient   *sql.DB
	ethNode    string
	ethAddress common.Address
)

const (
	ethNodeURL = "ETHURL"
)

func init() {

}

func main() {
	var (
		ok  bool
		err error
	)
	a := time.Now()
	pgClient, err = postgres.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to postgres", err)
	}
	defer pgClient.Close()
	// After connecting to the db lookup the latest block that we have
	lastStoredBlockNumber, err := postgres.GetLatestBlockInDB(pgClient)
	if err != nil {
		log.Fatal("Failed to get the latest block from the database")
	}

	if lastStoredBlockNumber == nil {
		var tmpInt int64 = 6988614 // Dec-31-2018 11:59:42 PM +UTC
		//var tmpInt int64 = 9282077 // Random block from earlyk January 2020
		lastStoredBlockNumber = &tmpInt
	}

	ethNode, ok = os.LookupEnv(ethNodeURL)
	if !ok {
		ethNode = "https://ropsten.infura.io/v3/<YOUR_API_KEY>"
	}

	ethClient, err := ethclient.Dial(ethNode)
	if err != nil {
		log.Fatal("Failed to connect to ethereum", err)
	}

	// Passing in nil returns the latest block.
	latestBlockNumberFromChain, err := blockchain.QueryLatestBlockHeader(ethClient)
	if err != nil {
		log.Fatalf("ERROR: Failed to get block by header: %v", err)
	}
	log.Printf("Latest Block Number [ %s ]", latestBlockNumberFromChain.Number.String())

	*lastStoredBlockNumber++ // We want to start at the next block.
	// Last block of 2019
	//lastBlockOf2019 := 9193265 // Dec-31-2019 11:59:45 PM +UTC
	//for *lastStoredBlockNumber < int64(lastBlockOf2019) {
	for *lastStoredBlockNumber < latestBlockNumberFromChain.Number.Int64() {
		block, err := blockchain.QueryBlockByNumber(ethClient, lastStoredBlockNumber)
		if err != nil {
			log.Printf("Could not query block number [ %d ] err: %v", lastStoredBlockNumber, err)
		}

		err = ethereum_history.ParseBlock(pgClient, ethClient, block)
		if err != nil {
			log.Printf("Failed to parse block number [ %d ] err: %v", lastStoredBlockNumber, err)
		}

		*lastStoredBlockNumber++

		if *lastStoredBlockNumber == latestBlockNumberFromChain.Number.Int64() {
			break
		}
	}
	delta := time.Now().Sub(a)
	log.Printf("Took  %v hours, %v minutes and %v seconds to parse all ethereum transactions for 2019",
		delta.Hours(), delta.Minutes(), delta.Seconds())
}
