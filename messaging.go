package main

import (
	"github.com/nats-io/go-nats"
)

const (
	chBitcoinAndroid  = "Bitcoin-Android"
	chEthereumAndroid = "Ethereum-Android"

	msgTitleEthereum = "Ethereum Confirmation"
	msgTitleBitcoin  = "Ethereum Confirmation"
)

type message struct {
	DeviceID string
	Message  string
	Title    string
}

// Publish a message to the nats channel Ethereum-Android
func publishAndroid(nc *nats.Conn, ec *nats.EncodedConn, channel string,
	msgTitle string, deviceID string, msg string) {

	sendChannel := make(chan *message)
	ec.BindSendChan(channel, sendChannel)

	me := &message{DeviceID: deviceID, Title: msgTitle, Message: msg}

	// Send via Go channels
	sendChannel <- me

	//receiveChannel := make(chan *person)
	//ec.BindRecvChan("hello", receivevChannel)
	// Receive via Go channels
	//who := <- receiveChannel
}
