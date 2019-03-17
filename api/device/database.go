package device

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

var collectionName = database.CollectionDevice

// CreateUserAccount add a user to the user table.
func UpdateDevice(client *mongo.Client, b *bytes.Buffer) error {
	var device Device
	collection := database.GetCollection(client, collectionName)
	err := json.Unmarshal(b.Bytes(), &device)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	res, err := collection.InsertOne(ctx, device)
	if err != nil {
		return err
	}

	log.Printf("New account inserted: %v", res.InsertedID)

	return nil
}

// GetDevices retrieve a list of devices for a given user.
func GetDevices(client *mongo.Client, id string) ([]Device, error) {
	devices := make([]Device, 0)
	collection := database.GetCollection(client, collectionName)
	user := bson.M{"_id": id}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, user)
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
		var account Device
		err := cur.Decode(&account)
		if err != nil {
			log.Printf("ERROR: failed to read cursor: %v", err)
		}
		devices = append(devices, account)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return devices, nil
}

func Delete(client *mongo.Client, id string) error {
	collection := database.GetCollection(client, collectionName)
	device := bson.M{"_id": id}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	res, err := collection.DeleteOne(ctx, device)
	if err != nil || res.DeletedCount != 1 {
		log.Printf("ERROR: failed to to remove device, maybe it didn't exist: %v", err)
	}

	return nil
}
