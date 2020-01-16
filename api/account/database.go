package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bahelit/confirmerator/database/mongodb"
)

const (
	collectionName = mongodb.CollectionAccount
)

// UpdateAccount add or update an account to the account table
func UpdateAccount(client *mongo.Client, b *bytes.Buffer) (string, error) {
	var account Account
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := mongodb.GetCollection(client, collectionName)
	err := json.Unmarshal(b.Bytes(), &account)
	if err != nil {
		log.Printf("Failed to parse: %v - err: %v", b, err)
		return "", err
	}

	if account.Symbol != nil {
		StringToChain(*account.Symbol)
		account.Symbol = nil
	}

	// Mobile QR-Code library prefixes addresses with detected wallet type.
	if strings.Contains(account.Address, ":") {
		tmpStr := strings.Split(account.Address, ":")
		account.Address = tmpStr[1]
	}

	// If id is populated then this is an update, else it's a new account
	if account.ID != nil && len(account.ID.String()) != 0 {
		filter := bson.D{{"_id", account.ID}}

		update := bson.D{
			{"$set", bson.D{
				{"userid", account.UserID},
				{"account_type", account.AccType},
				{"blockchain", account.Blockchain},
				{"address", account.Address},
				{"nickname", account.Nickname},
			}},
		}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Printf("ERROR: failed to update account: %v - err: %v", account, err)
			return "", err
		}

		log.Printf("Matched %v account and updated %v account.\n",
			updateResult.MatchedCount, updateResult.ModifiedCount)
	} else {
		res, err := collection.InsertOne(ctx, bson.D{
			{"userid", account.UserID},
			{"account_type", account.AccType},
			{"blockchain", account.Blockchain},
			{"address", account.Address},
			{"nickname", account.Nickname},
		})
		if err != nil {
			log.Printf("ERROR: failed to insert account: %v - err: %v", account, err)
			return "", err
		}

		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			idHex := map[string]interface{}{
				"id": oid.Hex(),
			}
			strOID := fmt.Sprintf("%s", idHex["id"])
			log.Printf("New account inserted: %v", strOID)
			return strOID, nil
		} else {
			log.Printf("New account inserted: %s", res.InsertedID)
		}
	}

	return "", nil
}

// GetAccountsForUser retrieve a list of accounts for a particular user.
// A user can have multiple accounts
func GetAccountsForUser(client *mongo.Client, userID string) ([]Account, error) {
	accounts := make([]Account, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := mongodb.GetCollection(client, collectionName)

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
		var account Account
		err := cur.Decode(&account)
		if err != nil {
			log.Printf("ERROR: failed to read cursor: %v", err)
		}
		log.Printf("Account: %v", account)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := mongodb.GetCollection(client, mongodb.CollectionAccount)

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
		//log.Println(account)
		accounts = append(accounts, account)
	}
	if err := cur.Err(); err != nil {
		log.Printf("ERROR: failed to read cursor for blockchain accounts - err: %v", err)
	}

	return accounts, nil
}
