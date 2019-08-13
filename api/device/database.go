package device

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

	"github.com/bahelit/confirmerator/database"
)

const (
	collectionName = database.CollectionDevice
)

// CreateUserAccount add a user to the user table.
func UpdateDevice(client *mongo.Client, b *bytes.Buffer) (string, error) {
	var device Device
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := database.GetCollection(client, collectionName)
	err := json.Unmarshal(b.Bytes(), &device)
	if err != nil {
		log.Printf("Failed to parse: %v - err: %v", b, err)
		return "", err
	}

	if device.ID != nil && len(device.ID.String()) != 0 {
		filter := bson.D{{"_id", device.ID}}

		update := bson.D{
			{"$set", bson.D{
				{"userid", device.UserID},
				{"platform", device.Platform},
				{"active", device.Active},
				{"identifier", device.Identifier},
			}},
		}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Printf("ERROR: failed to insert device: %v - err: %v", device, err)
			return "", err
		}

		log.Printf("Matched %v device and updated %v account.\n",
			updateResult.MatchedCount, updateResult.ModifiedCount)
	} else {
		res, err := collection.InsertOne(ctx, bson.D{
			{"userid", device.UserID},
			{"platform", device.Platform},
			{"active", device.Active},
			{"identifier", device.Identifier},
		})
		if err != nil {
			log.Printf("ERROR: failed to insert device: %v - err: %v", device, err)
			return "", err
		}

		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			idHex := map[string]interface{}{
				"id": oid.Hex(),
			}
			strOID := fmt.Sprintf("%s", idHex["id"])
			log.Printf("New device inserted: %v", strOID)
			return strOID, nil
		} else {
			log.Printf("New device inserted: %s", res.InsertedID)
		}
	}

	return "", nil
}

// GetDevice retrieve a single device associated with the user for a given platform.
func GetDevice(client *mongo.Client, platform int16, userID string) (string, error) {
	var device Device
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := database.GetCollection(client, database.CollectionDevice)

	filter := bson.M{"platform": platform, "userid": userID}
	err := collection.FindOne(ctx, filter).Decode(&device)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("ERROR: failed to query device: %v - err: %v", userID, err)
	} else if err != nil {
		log.Printf("no device found: %v - err: %v", userID, err)
	}

	return device.Identifier, nil
}

// GetDevices retrieve a list of devices for a given user.
func GetDevices(client *mongo.Client, userID string) ([]Device, error) {
	devices := make([]Device, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
		var account Device
		err := cur.Decode(&account)
		if err != nil {
			log.Printf("ERROR: failed to read cursor: %v", err)
		}
		devices = append(devices, account)
	}
	if err := cur.Err(); err != nil {
		log.Printf("ERROR: failed to read cursor for device - err: %v", err)
	}

	return devices, nil
}

func Delete(client *mongo.Client, id string) error {
	collection := database.GetCollection(client, collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	device := bson.M{"_id": id}

	res, err := collection.DeleteOne(ctx, device)
	if err != nil || res.DeletedCount != 1 {
		log.Printf("ERROR: failed to to remove device, maybe it didn't exist: %v", err)
	}

	return nil
}
