// Read information from the database, no writes done here.
package database

import (
	"database/sql"
	"fmt"
	"os"

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
	PlatformAndroid
)

type Account struct {
	RowID      int    `json:"rowID"`
	UserID     int    `json:"userID"`
	Address    string `json:"address"`
	Nickname   string `json:"nickname"`
	BlockChain int16  `json:"blockchain"`
	AccType    int16  `json:"accountType"`
	Device     string `json:"device"`
}

func InitDB(db *sql.DB) (*sql.DB, error) {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to postgres!")

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
func GetAccounts(db *sql.DB, blockchain int16) ([]Account, error) {
	accounts := make([]Account, 0)
	stmt, err := db.Prepare("SELECT id, user_id, address, type, nick_name FROM account WHERE blockchain = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(blockchain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		err = rows.Scan(&account.RowID, &account.UserID, &account.Address, &account.AccType, &account.Nickname)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// getAndroidDevices retrieve a list of android devices associated with the user.
func GetDevice(db *sql.DB, platform int16, userID int) (string, error) {
	var device string
	stmt, err := db.Prepare("SELECT DISTINCT ON (device.device_identifier) device_identifier	FROM device	WHERE active = TRUE AND platform = $1 AND user_id = $2")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	row := stmt.QueryRow(platform, userID)
	switch err := row.Scan(&device); err {
	case sql.ErrNoRows:
		fmt.Println("No devices found!")
	case nil:
		return device, nil
	default:
		return "", err
	}

	return device, nil
}
