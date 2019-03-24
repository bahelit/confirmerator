// Read information from the database, no writes done here.
package database

import (
	"context"
	"log"
	"os"
	"time"

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
	PlatformMobile
)

const (
	CollectionAccount = "account"
	CollectionDevice  = "device"
	CollectionUser    = "user"
)

const Database = "confirmerator"

func GetCollection(client *mongo.Client, collection string) *mongo.Collection {
	return client.Database(Database).Collection(collection)
}

// TODO Don't go down if we can't connect, implement retry.
func InitDB() (*mongo.Client, error) {
	config := dbConfig()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config[mongoURI]))

	// Check the connection
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("ERROR: failed to connect to mongodb: %v", err)
	}

	log.Println("Yay, Connected to MongoDB!")

	return client, nil
}

func dbConfig() map[string]string {
	conf := make(map[string]string)

	user, ok := os.LookupEnv(dbuser)
	if !ok {
		log.Print("DBUSER environment variable not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		log.Print("DBPASS environment not set")
	}
	mongoURI, ok := os.LookupEnv(mongoURI)
	if !ok {
		mongoURI = "mongodb://localhost:27017"
		log.Print("DBNAME environment variable set")
	}

	conf[dbuser] = user
	conf[dbpass] = password
	conf[mongoURI] = mongoURI
	return conf
}
