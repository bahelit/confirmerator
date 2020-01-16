package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bahelit/confirmerator/database/mongodb"
)

const (
	collectionName = mongodb.CollectionUser
)

// UpdateUserAccount add a user to the user table.
func UpdateUserAccount(client *mongo.Client, b *bytes.Buffer) (string, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := mongodb.GetCollection(client, collectionName)
	err := json.Unmarshal(b.Bytes(), &user)
	if err != nil {
		log.Printf("Failed to parse: %v - err: %v", b, err)
		return "", err
	}

	// If id is populated then this is an update, else it's a new account
	if user.ID != nil && len(user.ID.String()) != 0 {
		filter := bson.D{{"_id", user.ID}}

		update := bson.D{
			{"$set", bson.D{
				{"uid", user.UID},
				{"type", user.Type},
				{"email", user.Email},
				{"nickname", user.NickName},
			}},
		}

		updateResult, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("ERROR: failed to update user: %v - err: %v", user, err)
		}

		if updateResult.ModifiedCount == 0 {
			log.Printf("INFO: failed to find a match for: %v", user)
		}

		log.Printf("Matched %v user and updated %v user.\n",
			updateResult.MatchedCount, updateResult.ModifiedCount)
	} else {
		res, err := collection.InsertOne(ctx, bson.D{
			{"uid", user.UID},
			{"type", user.Type},
			{"email", user.Email},
			{"nickname", user.NickName},
		})
		if err != nil {
			log.Printf("ERROR: failed to insert user: %v - err: %v", user, err)
			return "", err
		}

		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			idHex := map[string]interface{}{
				"id": oid.Hex(),
			}
			strOID := fmt.Sprintf("%s", idHex["id"])
			log.Printf("New user inserted: %v", strOID)
			return strOID, nil
		} else {
			log.Printf("New user inserted: %s", res.InsertedID)
		}
	}

	return "", nil
}

// GetUserAccount get information about a user in the user.
func GetUserAccount(client *mongo.Client, uid string) (User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := mongodb.GetCollection(client, collectionName)

	err := collection.FindOne(ctx, bson.M{"uid": uid}).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("ERROR: failed to query user: %v - err: %v", uid, err)
	}

	return user, nil
}
