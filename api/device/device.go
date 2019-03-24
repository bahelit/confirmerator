package device

type Device struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"userID"`
	Platform   int16  `json:"platform"`
	Active     bool   `json:"active"`
	Identifier string `json:"identifier"`
}