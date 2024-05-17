package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/EclesioMeloJunior/btc-handshake/internal/messages"
	"github.com/EclesioMeloJunior/btc-handshake/internal/network"
)

func sendHandshakeAndWaitResponse() {
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	version := messages.NewVersion(
		messages.WithNumber(60002),
		messages.WithServices(messages.NodeNetwork|messages.NodeNetworkLimited),
		messages.WithAddrRecv(peerAddr, uint16(peerPort), 1),
		messages.WithAddrFrom("0.0.0.0", 8080, 1),
		messages.WithNonce(rng.Uint64()),
		messages.WithUserAgent("btc/eclesios-node"),
	)

	handshakeVersionMessage := messages.NewMainMessage([]byte("version"), version)
	encodedHandshake, err := handshakeVersionMessage.Encode()
	if err != nil {
		log.Fatalf("could not encode version handshake message: %s", err.Error())
	}

	srv, err := network.Dial(fmt.Sprintf("%s:%d", peerAddr, peerPort))
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

	// remote should send a version message back and a verack
	// as described here: https://en.bitcoin.it/wiki/Version_Handshake
	response := make([]byte, 1024)
	n, err := srv.WaitResponse(response)
	if err != nil {
		log.Fatalf("while waiting response: %s", err.Error())
	}

	responseReader := bytes.NewReader(response[:n])

	fmt.Printf("\nreceived from remote (%d bytes): 0x%x\n", n, response[:n])
	remoteVersionMessage := messages.NewMainMessage(nil, &messages.Version{})
	err = remoteVersionMessage.Decode(responseReader)
	if err != nil {
		log.Fatalf("while decoding remote's message: %s", err.Error())
	}
	fmt.Printf("remote's message:\n%s\n\n", remoteVersionMessage.String())

	remoteVerAck := messages.NewMainMessage(nil, messages.EmptyPayload{})
	err = remoteVerAck.Decode(responseReader)
	if err != nil {
		log.Fatalf("while decoding remote's message: %s", err.Error())
	}
	fmt.Printf("remote's message\n%s\n\n", remoteVerAck.String())

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
