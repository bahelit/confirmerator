package device

type Device struct {
	ID         int    `json:"id"`
	UserID     int    `json:"userID"`
	Platform   int16  `json:"platform"`
	Active     bool   `json:"active"`
	Identifier string `json:"identifier"`
}
