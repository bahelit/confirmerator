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

const (
	collectionName = database.CollectionUser
)

// UpdateUserAccount add a user to the user table.
func UpdateUserAccount(client *mongo.Client, b *bytes.Buffer) error {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := database.GetCollection(client, collectionName)
	err := json.Unmarshal(b.Bytes(), &user)
	if err != nil {
		return err
	}

	// If id is populated then this is an update, else it's a new account
	if len(user.ID) > 0 {
		filter := bson.D{{"_id", user.ID}}

		update := bson.D{
			{"$set", user},
		}

		updateResult, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("ERROR: failed to update user: %v - err: %v", user, err)
		}

		log.Printf("Matched %v user and updated %v user.\n",
			updateResult.MatchedCount, updateResult.ModifiedCount)
	} else {

		res, err := collection.InsertOne(ctx, user)
		if err != nil {
			log.Printf("ERROR: failed to insert user: %v - err: %v", user, err)
			return err
		}

		log.Printf("New user inserted: %v", res.InsertedID)
	}

	return nil
}

// GetUserAccount get information about a user in the user.
func GetUserAccount(client *mongo.Client, uid string) (User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := database.GetCollection(client, collectionName)

	err := collection.FindOne(ctx, bson.M{"uid": uid}).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("ERROR: failed to query user: %v - err: %v", uid, err)
	}

	return user, nil
}
