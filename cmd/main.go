package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/EclesioMeloJunior/btc-handshake/internal/messages"
	"github.com/EclesioMeloJunior/btc-handshake/internal/network"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	version := messages.NewVersion(
		messages.WithNumber(60002),
		messages.WithServices(messages.NodeNetwork|messages.NodeNetworkLimited),
		messages.WithAddrRecv("0.0.0.0", 8333, 1),
		messages.WithAddrFrom("0.0.0.0", 8080, 1),
		messages.WithNonce(rng.Uint64()),
		messages.WithUserAgent("btc/eclesio/eiger-handshake-challenger"),
	)

	handshakeVersionMessage := messages.NewMainMessage([]byte("version"), version)
	encodedHandshake, err := handshakeVersionMessage.Encode()
	if err != nil {
		log.Fatalf("could not encode version handshake message: %s", err.Error())
	}

	srv, err := network.New("0.0.0.0:8333")
	if err != nil {
		log.Fatalf("while instantiating network server: %s", err.Error())
	}

	// sending the encode handshake, as explained in the
	// protocol wiki, the remote side should send a version message back
	// as well as the verack message
	err = srv.Send(encodedHandshake)
	if err != nil {
		log.Fatalf("while sending version handshake: %s", err.Error())
	}

	// remote should send a version message back
	response := make([]byte, 1024)
	n, err := srv.WaitResponse(response)
	if err != nil {
		log.Fatalf("while waiting response: %s", err.Error())
	}

	fmt.Printf("received from remote (%d bytes): 0x%x\n", n, response[:n])
	remoteVersionMessage := &messages.Message{}

	// remote should send a verack
	n, err = srv.WaitResponse(response)
	if err != nil {
		log.Fatalf("while waiting response: %s", err.Error())
	}

	fmt.Printf("received from remote (%d bytes): 0x%x\n", n, response[:n])

	// we should send a verack since we received the remote's version
	verrackMessage := messages.NewMainMessage([]byte("verack"), messages.EmptyPayload{})
	encodeVerrack, err := verrackMessage.Encode()
	if err != nil {
		log.Fatalf("encoding verack: %s", err.Error())
	}

	err = srv.Send(encodeVerrack)
	if err != nil {
		log.Fatalf("while sending verack: %s", err.Error())
	}
}
