// Handle user related information
package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UID      string              `json:"uid"`
	Type     int16               `json:"type"`
	NickName string              `json:"nickname"`
}
