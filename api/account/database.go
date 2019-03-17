package account

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/bahelit/confirmerator/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionName = database.CollectionAccount

// UpdateAccount add or update an account to the account table
func UpdateAccount(client *mongo.Client, b *bytes.Buffer) error {
	var account Account
	collection := database.GetCollection(client, collectionName)
	err := json.Unmarshal(b.Bytes(), &account)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, account)
	if err != nil {
		return err
	}

	log.Printf("New account inserted: %v", res.InsertedID)

	return nil
}

// GetAccounts retrieve a list of accounts for a particular blockchain.
func GetAccounts(client *mongo.Client, userID string) ([]Account, error) {
	accounts := make([]Account, 0)
	collection := database.GetCollection(client, collectionName)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Printf("ERROR: failed to query accounts: %v", err)
	}
	defer func() {
		err := cur.Close(ctx)
		if err != nil {
			log.Printf("ERROR: failed to close cursor: %v", err)
		}
	}()

	for cur.Next(ctx) {
		//var result bson.M
		var account Account
		err := cur.Decode(&account)
		if err != nil {
			log.Printf("ERROR: failed to read cursor: %V", err)
		}
		accounts = append(accounts, account)
	}
	if err := cur.Err(); err != nil {
		log.Printf("ERROR: from cursor: %v", err)
	}

	return accounts, nil
}
