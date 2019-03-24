// Handle information related to the account, tracks wallet addresses
package chain_account

type Account struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"userID,omitempty"`
	AccType    int16  `json:"accountType"`
	Blockchain int16  `json:"blockchain"`
	Address    string `json:"address"`
	Nickname   string `json:"nickname"`
}
