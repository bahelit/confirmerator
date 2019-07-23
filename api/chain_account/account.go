// Handle information related to the account, tracks wallet addresses
package chain_account

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID         *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     string              `json:"userid"`
	AccType    int16               `json:"account_type"`
	Blockchain int16               `json:"blockchain"`
	Address    string              `json:"address"`
	Nickname   string              `json:"nickname"`
}
