package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bahelit/confirmerator/api/chain_account"
	"github.com/bahelit/confirmerator/database"
	"github.com/bahelit/confirmerator/shared"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nats-io/go-nats"
)

var (
	ethNode    string
	ethWSNode  string
	ethAddress common.Address

	natsServer string
)

const (
	ethNodeURL   = "ETHURL"
	ethWSNodeURL = "ETHWSURL"
	natsURL      = "NATSURL"
)

func main() {
	var ok bool
	client, err := database.InitDB()
	if err != nil || client == nil {
		log.Fatalf("Failed to connect to mongodb, bail'n %v", err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Printf("ERROR: failed to disconnect from mongo: %v", err)
		}
	}()

	ethNode, ok = os.LookupEnv(ethNodeURL)
	if !ok {
		ethNode = "https://ropsten.infura.io/v3/<YOUR_API_KEY>"
	}

	ethWSNode, ok = os.LookupEnv(ethWSNodeURL)
	if !ok {
		ethNode = "https://ropsten.infura.io/v3/<YOUR_API_KEY>"
	}

	natsServer, ok = os.LookupEnv(natsURL)
	if !ok {
		natsServer = nats.DefaultURL
	}
	natsConn, err := nats.Connect(natsServer, nats.ReconnectBufSize(5*1024*1024))
	if err != nil {
		log.Fatal("Failed to connect to gnats", err)
	}
	natsEncoder, _ := nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	defer natsEncoder.Close()

	ethClient, err := ethclient.Dial(ethNode)
	if err != nil {
		log.Fatal("Failed to connect to ethereum", err)
	}

	ethAccounts, err := chain_account.GetAccountsForBlockchain(client, database.ChainEthereum)
	if err != nil {
		// Can't do comparisons so just continue.
		log.Println(err)
	}
	fmt.Println("Number of Ethereum addresses to look up: ", len(ethAccounts))

	// On start-up print all the ethereum accounts we are tracking to the console
	for _, acct := range ethAccounts {
		if shared.IsValidAddress(ethAddress) {
			ethAddress := common.HexToAddress(acct.Address)
			ethValue := getBalance(ethClient, ethAddress)
			msg := fmt.Sprintln("Current balance: ", ethValue, " Address: ", acct.Nickname)
			fmt.Print(msg)
			//publishEthereumAndroid(natsConn, natsEncoder, testDevice, msg)
			//time.Sleep(30000)
		}
	}

	err = wsSubscribe(client, natsConn, natsEncoder, ethClient)
	if err != nil {
		log.Fatal(err)
	}
}
