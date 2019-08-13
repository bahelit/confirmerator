package main

import (
	"log"
	"os"

	"github.com/bahelit/confirmerator/messaging/firebase"
	"github.com/nats-io/nats.go"

	"github.com/bahelit/confirmerator/messaging"
)

var (
	natsServer string
)

const (
	natsURL = "NATSURL"
)

func main() {
	var ok bool

	// NATS Connect Options.
	opts := []nats.Option{nats.Name("NATS Confirmerator Subscriber"), nats.ReconnectBufSize(5 * 1024 * 1024)}
	opts = messaging.SetupConnOptions(opts)

	natsServer, ok = os.LookupEnv(natsURL)
	if !ok {
		natsServer = nats.DefaultURL
	}
	natsConn, err := nats.Connect(natsServer, opts...)
	if err != nil {
		log.Fatal("ERROR: Failed to connect to nats-server", err)
	}
	natsEncoder, _ := nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	defer natsEncoder.Close()

	firebase.SubscribeToChannel(natsEncoder, firebase.SbjEthereumAndroid)
}
