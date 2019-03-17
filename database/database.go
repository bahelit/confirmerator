// Read information from the database, no writes done here.
package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	dbuser   = "DBUSER"
	dbpass   = "DBPASS"
	mongoURI = "MONGOURI"
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

const (
	CollectionAccount = "account"
	CollectionDevice  = "device"
	CollectionUser    = "user"
)

type Account struct {
	ID         int    `json:"id"`
	UserID     int    `json:"userID"`
	Address    string `json:"address"`
	Nickname   string `json:"nickname"`
	Blockchain int16  `json:"blockchain"`
	AccType    int16  `json:"accountType"`
	Device     string `json:"device"`
}

var (
	Database = "confirmerator"
)

func GetCollection(client *mongo.Client, collection string) *mongo.Collection {
	return client.Database(Database).Collection(collection)
}

func InitDB(client *mongo.Client) error {
	config := dbConfig()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config[mongoURI]))

	// Check the connection
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("ERROR: failed to connect to mongodb: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	return nil
}

func dbConfig() map[string]string {
	conf := make(map[string]string)

	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	mongoURI, ok := os.LookupEnv(mongoURI)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}

	conf[dbuser] = user
	conf[dbpass] = password
	conf[mongoURI] = mongoURI
	return conf
}

// GetAccounts retrieve a list of accounts for a particular blockchain.
func GetBlockchainAccounts(client *mongo.Client, blockchain int16) ([]Account, error) {
	accounts := make([]Account, 0)

	collection := GetCollection(client, CollectionAccount)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{"blockchain": blockchain})
	if err != nil {
		log.Fatalf("ERROR: failed to get collection: %v", err)
	}
	defer func() {
		err := cur.Close(ctx)
		if err != nil {
			log.Printf("ERROR: failed to close cursor: %v", err)
		}
	}()

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return accounts, nil
}

// GetDevice retrieve a single device associated with the user for a given platform.
func GetDevice(client *mongo.Client, platform int16, userID int) (string, error) {
	var device string
	collection := GetCollection(client, CollectionDevice)

	var result struct {
		Value float64
	}
	filter := bson.M{"name": "pi"}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return device, nil
}
