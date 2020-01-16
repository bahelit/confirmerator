package ethereum

import (
	"context"
	"fmt"
	"github.com/bahelit/confirmerator/blockchain"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bahelit/confirmerator/api/account"
	"github.com/bahelit/confirmerator/api/device"
	"github.com/bahelit/confirmerator/database/mongodb"
	"github.com/bahelit/confirmerator/messaging/firebase"
	"github.com/bahelit/confirmerator/shared"
)

func queryTransactions(db *mongo.Client, client *ethclient.Client, block *types.Block, accounts []account.Account,
	ec *nats.EncodedConn) error {
	// We can read the transactions in a block by calling the Transaction method which returns a list of Transaction type.
	// It's then trivial to iterate over the collection and retrieve any information regarding the transaction.
	fmt.Println("Transactions in block: ", block.Transactions().Len())
	for _, tx := range block.Transactions() {
		//fmt.Println(tx.Hash().Hex())                // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		//fmt.Println(tx.Value().String()) // 10000000000000000
		//fmt.Println(tx.Gas())                       // 105000
		//fmt.Println(tx.GasPrice().Uint64())         // 102000000000
		//fmt.Println(tx.Nonce())                     // 110644
		// If it contains data it is a smart contract.
		//fmt.Println("Transaction data: ", tx.Data()) // []
		//if tx.To() != nil {
		//	log.Println("To address: ", tx.To().Hex()) // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
		//}
		if tx.To() == nil {
			log.Println("tx to address is empty.")
			continue
		}

		// In order to read the sender address, we need to call AsMessage on the transaction which returns a
		// Message type containing a function to return the sender (from) address.
		//if msg, err := tx.AsMessage(types.HomesteadSigner{}); err != nil {
		//	log.Printf("ERROR: Failed to read from AsMessage error: %v", err) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		//} else {
		//	log.Println("From address", msg.From().Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		//}

		//// NOTE this calls consumes allot of data from infura
		////
		// Each transaction has a receipt which contains the result of the execution of the transaction,
		// such as any return values and logs, as well as the status which will be 1 (success) or 0 (fail).
		//receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		//if err != nil {
		//	log.Printf("ERROR: Failed to get transaction receipt error: %v", err)
		//	continue
		//}
		//fmt.Println("Receipt Status: ", receipt.Status)

		for _, acct := range accounts {
			//log.Printf("Comparing Address: %v", acct.Address)
			if tx.To().String() == acct.Address {
				//ethAddress := shared.HexToAddress(acct.address)
				log.Printf("Confirmed Transaction for: %v", acct.Address)
				log.Printf("Looking for user: %v", acct.UserID)
				deviceIdentifier, err := device.GetDevice(db, mongodb.PlatformMobile, acct.UserID)
				if err != nil {
					log.Printf("failed to find device for userID [ %s ] error: %s", acct.UserID, err)
					continue
				}

				// Hide the actual identifier from the logs, show the first 4 and last 4.
				log.Printf("found deviceID: %s", deviceIdentifier[:4]+"..."+deviceIdentifier[len(deviceIdentifier)-4:])
				//msg := fmt.Sprintf("Confirmed wallet %s txHash %s", acct.Nickname, tx.Hash().String())
				msg := fmt.Sprintf("Confirmed transaction for %s", acct.Nickname)

				weiBalance := new(big.Float)
				weiBalance.SetString(tx.Value().String())
				ethValue := new(big.Float).Quo(weiBalance, big.NewFloat(math.Pow10(18)))

				firebase.PublishFirebase(ec, firebase.SbjEthereumAndroid, shared.SymbEthereum, firebase.MsgTitleEthereum, deviceIdentifier, acct.Nickname, msg, ethValue.String())
			}
		}
	}

	return nil
}

func SubscribeWebSocket(db *mongo.Client, ec *nats.EncodedConn, wsClient *ethclient.Client) error {
	headers := make(chan *types.Header)

	sub, err := wsClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Printf("ERROR: No connection to web socket error: %v", err)
		return err
	}

	for {
		select {
		case err := <-sub.Err():
			log.Println(err)
			return err
		case header := <-headers:
			//log.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			log.Println("Block number: ", header.Number)
			time.Sleep(time.Second)

			block, err := blockchain.QueryBlockByHeader(wsClient, header)
			if err != nil {
				// Cluttering the logs
				log.Printf("ERROR: failed to query block: %v", err)
				continue
			}

			if len(block.Transactions()) > 0 {
				ethAccounts, err := account.GetAccountsForBlockchain(db, account.ChainEthereum)
				if err != nil {
					// Can't do comparisons so just continue.
					log.Printf("ERROR: Failed to query for eth account error: %v", err)
					continue
				}

				err = queryTransactions(db, wsClient, block, ethAccounts, ec)
				if err != nil {
					log.Printf("ERROR: Failed to queury trasaction on block [ %v ] error: %v", block.Number(), err)
				}
			} else {
				fmt.Println("Zero transactions in block")
			}
		}
	}
}

func WebSocketReconnect(ethWSNode string) (*ethclient.Client, error) {
	wsClient, err := ethclient.Dial(ethWSNode)
	if err != nil {
		log.Printf("ERROR: Failed to connect to ethereum web socket error: %v", err)
		return nil, err
	}

	return wsClient, err
}
