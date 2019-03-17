// Handle information related to the account, tracks wallet addresses
package account

type Account struct {
	ID         int    `json:"id"`
	UserID     int    `json:"userID,omitempty"`
	AccType    int16  `json:"accountType"`
	Blockchain int16  `json:"blockchain"`
	Address    string `json:"address"`
	Nickname   string `json:"nickname"`
}
