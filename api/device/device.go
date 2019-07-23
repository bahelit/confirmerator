package device

import "go.mongodb.org/mongo-driver/bson/primitive"

type Device struct {
	ID         *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     string              `json:"userid"`
	Platform   int16               `json:"platform"`
	Active     bool                `json:"active"`
	Identifier string              `json:"identifier"`
}
