package firebase

import (
	"github.com/bahelit/confirmerator/shared"
	"log"

	"github.com/nats-io/nats.go"
)

const (
	SbjBitcoinAndroid  = "Bitcoin-Android"
	SbjEthereumAndroid = "Ethereum-Android"

	MsgTitleEthereum = "Ethereum Confirmation"
	MsgTitleBitcoin  = "Ethereum Confirmation"
)

// Publish a message to Firebase ids
func PublishFirebase(ec *nats.EncodedConn, subject string, symbol shared.Symbol,
	msgTitle string, deviceID string, nickname string, msg string, value string) {

	me := &Message{DeviceID: deviceID, Title: msgTitle, Message: msg, Symbol: symbol, Value: value, Nickname: nickname}
	log.Printf("sending message to : %v", me.DeviceID[:4]+"..."+me.DeviceID[len(me.DeviceID)-4:])

	sendChannel := make(chan *Message)
	err := ec.BindSendChan(subject, sendChannel)
	if err != nil {
		log.Printf("failed to bind to channel error: %v", err)
		return
	}

	//// Uncomment to push directly without queue
	//go PushMessage(me)

	// Send via Go channels
	sendChannel <- me
}
