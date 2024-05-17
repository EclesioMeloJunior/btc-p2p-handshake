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

func listenForHandshakes() {
	rcv, err := network.Listen()
	if err != nil {
		log.Fatalf(err.Error())
	}

	for v := range rcv {
		incoming := make([]byte, 1024)
		n, err := v.WaitResponse(incoming)
		if err != nil {
			log.Fatalf("while reading connection: %s", err.Error())
		}

		fmt.Printf("\nreceived from remote (%d bytes): 0x%x\n", n, incoming[:n])
		remoteVersionMessage := messages.NewMainMessage(nil, &messages.Version{})
		err = remoteVersionMessage.Decode(bytes.NewReader(incoming[:n]))
		if err != nil {
			log.Fatalf("while decoding remote's message: %s", err.Error())
		}
		fmt.Printf("remote's message:\n%s\n\n", remoteVersionMessage.String())

		err = sendOurVersionAndVerAck(v)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func sendOurVersionAndVerAck(stream network.Stream) error {
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	version := messages.NewVersion(
		messages.WithNumber(60002),
		messages.WithServices(messages.NodeNetwork|messages.NodeNetworkLimited),
		messages.WithAddrRecvFromString(stream.RemoteAddr().String(), 0),
		messages.WithAddrFrom("0.0.0.0", 8080, 1),
		messages.WithNonce(rng.Uint64()),
		messages.WithUserAgent("btc/eclesios-node"),
	)

	handshakeVersionMessage := messages.NewMainMessage([]byte("version"), version)
	encodedHandshake, err := handshakeVersionMessage.Encode()
	if err != nil {
		log.Fatalf("could not encode version handshake message: %s", err.Error())
	}

	err = stream.Send(encodedHandshake)
	if err != nil {
		return fmt.Errorf("while sending our version: %w", err)
	}

	// we should send a verack since we received the remote's version
	verrackMessage := messages.NewMainMessage([]byte("verack"), messages.EmptyPayload{})
	encodeVerrack, err := verrackMessage.Encode()
	if err != nil {
		log.Fatalf("encoding verack: %s", err.Error())
	}

	err = stream.Send(encodeVerrack)
	if err != nil {
		log.Fatalf("while sending verack: %s", err.Error())
	}

	return nil
}
