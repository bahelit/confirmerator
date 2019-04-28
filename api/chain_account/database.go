package chain_account

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

const (
	collectionName = database.CollectionAccount
)

// UpdateAccount add or update an account to the account table
func UpdateAccount(client *mongo.Client, b *bytes.Buffer) error {
	var account Account
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := database.GetCollection(client, collectionName)
	err := json.Unmarshal(b.Bytes(), &account)
	if err != nil {
		log.Printf("Failed to parse: %v - err: %v", b, err)
		return err
	}

	// If id is populated then this is an update, else it's a new account
	if len(account.ID) > 0 {
		filter := bson.D{{"_id", account.ID}}

		update := bson.D{
			{"$set", account},
		}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Printf("ERROR: failed to update account: %v - err: %v", account, err)
		}

		log.Printf("Matched %v account and updated %v account.\n",
			updateResult.MatchedCount, updateResult.ModifiedCount)
	} else {
		res, err := collection.InsertOne(ctx, account)
		if err != nil {
			log.Printf("ERROR: failed to insert account: %v - err: %v", account, err)
			return err
		}

		log.Printf("New account inserted: %v", res.InsertedID)
	}

	return nil
}

// GetAccountsForUser retrieve a list of accounts for a particular user.
// A user can have multiple accounts
func GetAccountsForUser(client *mongo.Client, userID string) ([]Account, error) {
	accounts := make([]Account, 0)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := database.GetCollection(client, collectionName)

	cur, err := collection.Find(ctx, bson.M{"userid": userID})
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

// GetAccountsForBlockchain retrieve a list of accounts for a particular blockchain.
func GetAccountsForBlockchain(client *mongo.Client, blockchain int16) ([]Account, error) {
	accounts := make([]Account, 0)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := database.GetCollection(client, database.CollectionAccount)

	cur, err := collection.Find(ctx, bson.M{"blockchain": blockchain})
	if err != nil {
		log.Printf("ERROR: failed to get blockchain collection: %v", err)
	}
	defer func() {
		err := cur.Close(ctx)
		if err != nil {
			log.Printf("ERROR: failed to close cursor: %v", err)
		}
	}()

	for cur.Next(ctx) {
		var account Account
		err := cur.Decode(&account)
		if err != nil {
			log.Printf("ERROR: failed to decode account form cursor: %v", err)
		}
		log.Println(account)
		accounts = append(accounts, account)
	}
	if err := cur.Err(); err != nil {
		log.Printf("ERROR: failed to read cursor for blockchain accounts - err: %v", err)
	}

	return accounts, nil
}
