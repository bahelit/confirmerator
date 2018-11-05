package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/bahelit/confirmerator/database"
	"github.com/bahelit/confirmerator/shared"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	"github.com/nats-io/go-nats"
)

var (
	dbHandle   *sql.DB
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
	dbHandle, err := database.InitDB(dbHandle)
	if err != nil {
		log.Fatal("Failed to connect to postgres", err)
	}
	defer dbHandle.Close()

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

	client, err := ethclient.Dial(ethNode)
	if err != nil {
		log.Fatal("Failed to connect to ethereum", err)
	}

	ethAccounts, err := database.GetAccounts(dbHandle, database.ChainEthereum)
	if err != nil {
		// Can't do comparisons so just continue.
		log.Println(err)
	}
	fmt.Println("Number of Ethereum addresses to look up: ", len(ethAccounts))

	// On start-up print all the ethereum accounts we are tracking tot the console
	for _, acct := range ethAccounts {
		if shared.IsValidAddress(ethAddress) {
			ethAddress := common.HexToAddress(acct.Address)
			ethValue := getBalance(client, ethAddress)
			msg := fmt.Sprintln("Current balance: ", ethValue, " Address: ", acct.Nickname)
			fmt.Print(msg)
			//publishEthereumAndroid(natsConn, natsEncoder, testDevice, msg)
			//time.Sleep(30000)
		}
	}

	err = wsSubscribe(dbHandle, natsConn, natsEncoder, client)
	if err != nil {
		log.Fatal(err)
	}
}
