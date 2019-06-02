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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config[mongoURI]))
	if err != nil || client == nil {
		log.Printf("ERROR: failed to connect to mongodb: %v", err)
	}

	// Check the connection
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("ERROR: failed to connect to mongodb: %v", err)
	}

	log.Println("Yay, Connected to MongoDB!")

	return client, nil
}

func dbConfig() map[string]string {
	conf := make(map[string]string)

	mongoURI, ok := os.LookupEnv(mongoURI)
	if !ok {
		mongoURI = "mongodb://localhost:27017"
		log.Print("MONGOURI environment variable set")
	}

	conf[mongoURI] = mongoURI
	return conf
}
