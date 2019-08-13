// Send message to mobile devices through fcm
package firebase

import (
	"log"
	"os"

	"gopkg.in/maddevsio/fcm.v1"

	"github.com/bahelit/confirmerator/shared"
)

var (
	firebaseKey = "FIREBASE_KEY"
	apiKey      string
)

type Message struct {
	DeviceID string
	Nickname string
	Message  string
	Symbol   shared.Symbol
	Title    string
	Value    string
}

func init() {
	var ok bool
	apiKey, ok = os.LookupEnv(firebaseKey)
	if !ok {
		log.Printf("failed to retrieve firebase api key, check env")
	}
}

func PushMessage(msg *Message) {
	data := map[string]string{
		"msg":          "Derpy Derp",
		"id":           msg.Nickname,
		"value":        msg.Value,
		"symbol":       msg.Symbol.ToString(),
		"click_action": "FLUTTER_NOTIFICATION_CLICK",
	}
	c := fcm.NewFCM(apiKey)
	response, err := c.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  []string{msg.DeviceID},
		ContentAvailable: true,
		Priority:         fcm.PriorityNormal,
		Notification: fcm.Notification{
			Title:       msg.Title,
			Body:        msg.Message,
			ClickAction: "FLUTTER_NOTIFICATION_CLICK",
		},
	})
	if err != nil {
		log.Fatalf("ERROR: failed to deliver message error: %v", err)
	}
	log.Println("Status Code   :", response.StatusCode)
	log.Println("Success       :", response.Success)
	log.Println("Fail          :", response.Fail)
	log.Println("Canonical_ids :", response.CanonicalIDs)
	log.Println("Topic MsgId   :", response.MsgID)
}
