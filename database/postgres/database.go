package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

const (
	ChainBitcoin = iota
	ChainEthereum
	ChainEthereumClassic
	ChainBitcoinCash
	ChainCallisto
	ChainRavenCoin
)

const (
	PlatformWeb = iota + 1
	PlatformDesktop
)

type EthTransaction struct {
	Time          int64
	TxnFrom       common.Address
	TxnTo         common.Address
	Gas           uint64
	GasPrice      *big.Int
	Block         int64
	TxnHash       string
	Value         float64
	ContractTo    common.Address `json:"contract_to"`
	ContractValue string         `json:"contract_value"`
}

func InitDB() (*sql.DB, error) {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to postgres!")

	return db, nil
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

// getAccounts retrieve a list of accounts for a particular blockchain.
func GetLatestBlockInDB(db *sql.DB) (latestBlock *int64, err error) {
	stmt, err := db.Prepare("SELECT Max(block) from public.ethtxns")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	switch err = row.Scan(&latestBlock); err {
	case sql.ErrNoRows:
		log.Println("No blocks found!")
		return nil, nil
	case nil:
		return latestBlock, nil
	default:
		return nil, err
	}
}

func BulkInsertTransactions(db *sql.DB, ethTxns []EthTransaction) error {
	a := time.Now()
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, _ := txn.Prepare(pq.CopyIn("ethtxns", "time", "txnfrom", "txnto", "gas", "gasprice", "block", "txnhash", "value", "contract_to", "contract_value")) // MessageDetailRecord is the table name

	for txn := range ethTxns {
		//log.Printf("TrannYYY: %v", ethTxns[txn].Value)
		var transaction EthTransaction
		transaction = ethTxns[txn]

		_, err := stmt.Exec(transaction.Time, transaction.TxnFrom.String(), transaction.TxnTo.String(),
			transaction.Gas, transaction.GasPrice.Int64(), transaction.Block,
			transaction.TxnHash, transaction.Value, transaction.ContractTo.String(),
			strings.Replace(transaction.ContractValue, "\u0000", "", -1))
		if err != nil {
			log.Printf("Failed to exec stmt in loop, err: %v", err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Failed to exec stmt, err: %v", err)
		for txnErrList := range ethTxns {
			log.Printf("to %v \nfrom %v \nvalue %v \nhash %v \ncontractTo %v \ncontractValue %v ",
				ethTxns[txnErrList].TxnTo.String(), ethTxns[txnErrList].TxnFrom.String(), ethTxns[txnErrList].Value, ethTxns[txnErrList].TxnHash,
				ethTxns[txnErrList].ContractTo.String(), ethTxns[txnErrList].ContractValue)
		}
	}
	err = stmt.Close()
	if err != nil {
		log.Printf("Failed to close stmt, err: %v", err)
	}
	err = txn.Commit()
	if err != nil {
		log.Printf("Failed to commit stmt, err: %v", err)
	}

	delta := time.Now().Sub(a)
	log.Printf("Took [ %d ] milliseconds to insert [ %d ] records", delta.Milliseconds(), len(ethTxns))

	return nil
}

// getAccounts retrieve a list of accounts for a particular blockchain.
func GetEthTransactions(db *sql.DB, address common.Address) ([]EthTransaction, error) {
	transactions := make([]EthTransaction, 0)
	stmt, err := db.Prepare("Select time, txnfrom, txnto, gas, gasprice, block, txnhash, value, contract_to, contract_value FROM ethtxns WHERE txnto OR  contract_to = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(address.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction EthTransaction
		err = rows.Scan(&transaction.Time, &transaction.TxnFrom, &transaction.TxnTo, &transaction.Gas, &transaction.GasPrice,
			&transaction.Block, &transaction.TxnHash, &transaction.Value, &transaction.ContractTo, &transaction.ContractValue)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
