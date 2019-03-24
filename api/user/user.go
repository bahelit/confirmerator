// Handle user related information
package user

type User struct {
	ID       string `json:"id" bson:"_id"`
	UID      string `json:"uid"`
	Type     int16  `json:"type"`
	NickName string `json:"nickname"`
}
