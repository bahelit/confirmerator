// Handle user related information
package user

type User struct {
	ID       int    `json:"id"`
	UID      string `json:"uid"`
	Type     int16  `json:"type"`
	NickName string `json:"nickname"`
}
