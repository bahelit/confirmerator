package main

import (
	"context"
	"fmt"
	"github.com/bahelit/confirmerator/cmd/confirmerator-blockchain-subscriber/ethereum"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nats-io/nats.go"

	"github.com/bahelit/confirmerator/api/account"
	"github.com/bahelit/confirmerator/blockchain"
	"github.com/bahelit/confirmerator/database/mongodb"
	"github.com/bahelit/confirmerator/messaging"
	"github.com/bahelit/confirmerator/shared"
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

func init() {
	var ok bool
	ethNode, ok = os.LookupEnv(ethNodeURL)
	if !ok {
		ethNode = "https://ropsten.infura.io/v3/<YOUR_API_KEY>"
	}

	ethWSNode, ok = os.LookupEnv(ethWSNodeURL)
	if !ok {
		ethNode = "wss://ropsten.infura.io/ws/v3/<YOUR_API_KEY>"
	}

	natsServer, ok = os.LookupEnv(natsURL)
	if !ok {
		natsServer = nats.DefaultURL
	}
}

func main() {
	client, err := mongodb.InitDB()
	if err != nil || client == nil {
		log.Fatalf("ERROR: Failed to connect to mongodb, bail'n %v", err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Printf("ERROR: failed to disconnect from mongo: %v", err)
		}
	}()

	// NATS Connect Options.
	opts := []nats.Option{nats.Name("NATS Confirmerator Subscriber"), nats.ReconnectBufSize(5 * 1024 * 1024)}
	opts = messaging.SetupConnOptions(opts)

	natsConn, err := nats.Connect(natsServer, opts...)
	if err != nil {
		log.Fatal("ERROR: Failed to connect to nats-server", err)
	}
	natsEncoder, _ := nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	defer natsEncoder.Close()

	//ethClient, err := ethclient.Dial(ethNode)
	//if err != nil {
	//	log.Fatal("ERROR: Failed to connect to ethereum", err)
	//}

	ethAccounts, err := account.GetAccountsForBlockchain(client, account.ChainEthereum)
	if err != nil {
		// Can't do comparisons so just continue.
		log.Println(err)
	}
	fmt.Println("Number of Ethereum addresses to look up: ", len(ethAccounts))

	// TODO implement better retry
	wsClient, err := ethclient.Dial(ethWSNode)
	if err != nil {
		log.Printf("ERROR: Failed to connect to ethereum web socket: %s", err)

		time.Sleep(30 * time.Second)

		wsClient, err = ethereum.WebSocketReconnect(ethWSNode)
		if err != nil {
			log.Fatalf("ERROR: Failed to connect to ethereum node error: %s", err)
		}
	}

	// On start-up print all the ethereum accounts we are tracking to the console
	for _, acct := range ethAccounts {
		if shared.IsValidAddress(ethAddress) {
			ethAddress := common.HexToAddress(acct.Address)
			ethValue := blockchain.GetBalance(wsClient, ethAddress)
			log.Printf("Current Balance: [ %v ] Wallet Nickname: %s", ethValue, acct.Nickname)
			//publishEthereumAndroid(natsConn, natsEncoder, testDevice, msg)
			//time.Sleep(30000)
		}
	}

	err = ethereum.SubscribeWebSocket(client, natsEncoder, wsClient)
	if err != nil {
		log.Fatalf("ERROR: Failed to subscribe to websocket error: %s", err)
	}
}
