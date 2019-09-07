package firebase

import (
	"log"

	"github.com/nats-io/nats.go"
)

func SubscribeToChannel(ec *nats.EncodedConn, subject string) {
	receiveChannel := make(chan *Message)
	sub, err := ec.BindRecvChan(subject, receiveChannel)
	if err != nil {
		log.Printf("failed to bind to channel error: %v", err)
		return
	}
	defer sub.Unsubscribe()

	// Receive via Go channels
	notification := <-receiveChannel

	handleMessage(notification)

	//// Go type Subscriber
	//sub, err := ec.Subscribe(subject, handleMessage)
	//if err != nil {
	//	log.Printf("failed to bind to channel error: %v", err)
	//	return
	//}
	//defer sub.Unsubscribe()
}

func handleMessage(msg *Message) {
	log.Printf("Received a Message: %+v\n", msg)
	PushMessage(msg)
}
