package user

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

var collectionName = database.CollectionUser

// UpdateUserAccount add a user to the user table.
func UpdateUserAccount(client *mongo.Client, b *bytes.Buffer) error {
	var user User
	collection := database.GetCollection(client, collectionName)
	err := json.Unmarshal(b.Bytes(), &user)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	log.Printf("New account inserted: %v", res.InsertedID)

	return nil
}

// GetUserAccount get information about a user in the user.
func GetUserAccount(client *mongo.Client, uid string) (User, error) {
	var user User
	collection := database.GetCollection(client, collectionName)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, bson.M{"uid": uid}).Decode(&user)
	if err != nil {
		log.Printf("ERROR: failed to query accounts: %v", err)
	}

	return user, nil
}
