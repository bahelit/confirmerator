package ethereum_history

import (
	"context"
	"database/sql"
	"github.com/bahelit/confirmerator/shared"
	_ "github.com/lib/pq"
	"log"
	"math"
	"math/big"
	"regexp"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/bahelit/confirmerator/database/postgres"
)

func ParseTransaction(client *ethclient.Client, ethTxnChan chan postgres.EthTransaction, tx *types.Transaction, blockNumber, blockTime int64, wg *sync.WaitGroup) {
	defer wg.Done()

	transaction := postgres.EthTransaction{Block: blockNumber, Time: blockTime}
	if tx.To() != nil {
		transaction.TxnTo = *tx.To()
	}
	transaction.Gas = tx.Gas()
	transaction.GasPrice = tx.GasPrice()
	transaction.TxnHash = tx.Hash().String()

	weiBalance := new(big.Float)
	weiBalance.SetString(tx.Value().String())
	ethValue := new(big.Float).Quo(weiBalance, big.NewFloat(math.Pow10(18)))
	transaction.Value, _ = ethValue.Float64()

	// In order to read the sender address, we need to call AsMessage on the transaction which returns a
	// Message type containing a function to return the sender (from) address.
	if msg, err := tx.AsMessage(types.HomesteadSigner{}); err != nil {
		//log.Printf("ERROR: Failed to read from AsMessage to get TxnFrom, error: %v", err) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
	} else {
		//log.Println("From address", msg.From().Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		if shared.IsValidAddress(msg.From()) {
			transaction.TxnFrom = msg.From()
		}
	}

	//// NOTE this calls consumes allot of data from the endpoint!
	////
	// Each transaction has a receipt which contains the result of the execution of the transaction,
	// such as any return values and logs, as well as the status which will be 1 (success) or 0 (fail).
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil || receipt == nil {
		log.Printf("ERROR: Failed to get transaction receipt error: %v", err)
		ethTxnChan <- transaction
		return
	}

	// We can get info about token transfers from the receipt logs.
	// Topics[0] = Hash of the transaction
	// Topics[1] = From Address
	// Topics[2] = To Address
	if receipt.Logs != nil && len(receipt.Logs) != 0 {
		//log.Printf("Receipt Log Data: %v", receipt.Logs[0].Data)

		re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
		if re.MatchString(string(receipt.Logs[0].Data)) {
			transaction.ContractValue = string(receipt.Logs[0].Data)
		}

		if receipt.Logs[0].Topics != nil && len(receipt.Logs[0].Topics) == 3 {
			//log.Printf("Receipt Log Receipt: %v", receipt.Logs[0].Topics[2].String())
			contractTo := receipt.Logs[0].Topics[2].Hex()
			if shared.IsValidAddress(contractTo) {
				transaction.ContractTo.SetBytes([]byte(contractTo))
			}
		}
	}

	ethTxnChan <- transaction
	return
}

func ParseBlock(db *sql.DB, client *ethclient.Client, block *types.Block) error {
	// We can read the transactionsChannel in a block by calling the Transaction method which returns a list of Transaction type.
	// It's then trivial to iterate over the collection and retrieve any information regarding the transaction.
	log.Printf("[ %d ] transactions in block number [ %d ]", block.Transactions().Len(), block.Number().Int64())

	// Create a buffered channel to build up for a bulk insert.
	var wg sync.WaitGroup
	wg.Add(len(block.Transactions()))
	transactionsChannel := make(chan postgres.EthTransaction, len(block.Transactions()))

	time := block.Time()
	log.Printf("Block TimeStamp: %v", time)

	for tx := range block.Transactions() {
		go ParseTransaction(client, transactionsChannel, block.Transactions()[tx], block.Number().Int64(), int64(time), &wg)
	}

	wg.Wait()
	log.Println("...Done waiting")
	close(transactionsChannel)

	var transactions []postgres.EthTransaction
	for txn := range transactionsChannel {
		var transaction postgres.EthTransaction
		transaction = txn
		transactions = append(transactions, transaction)
	}

	err := postgres.BulkInsertTransactions(db, transactions)
	if err != nil {
		log.Printf("Failed to insert transactions, err: %v", err)
	}

	return nil
}
