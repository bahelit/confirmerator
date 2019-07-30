// Handle information related to the account, tracks wallet addresses
package account

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID         *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     string              `json:"userid"`
	AccType    int16               `json:"account_type"`
	Blockchain int16               `json:"blockchain"`
	Symbol     *string             `json:"symbol,omitempty" bson:"symbol,omitempty"`
	Address    string              `json:"address"`
	Nickname   string              `json:"nickname"`
}
