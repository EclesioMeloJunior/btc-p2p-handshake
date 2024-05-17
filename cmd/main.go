package main

import (
	"flag"
	"strings"
)

var (
	peerAddr string
	peerPort uint
)

func init() {
	flag.StringVar(&peerAddr, "peer-addr", "", "address in the format 0.0.0.0")
	flag.UintVar(&peerPort, "peer-port", 0, "peer valid TCP port")
}

func main() {
	flag.Parse()

	if strings.TrimSpace(peerAddr) != "" && peerPort != 0 {
		sendHandshakeAndWaitResponse()
	} else {
		listenForHandshakes()
	}
}
