package network

import (
	"fmt"
	"net"
	"net/netip"
)

type Stream struct {
	remote  net.Addr
	tcpConn net.Conn
}

func (s Stream) RemoteAddr() net.Addr {
	return s.remote
}

func Listen() (<-chan Stream, error) {
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, fmt.Errorf("while setup tcp listener: %w", err)
	}

	fmt.Println("listening on 0.0.0.0:8080")

	connCh := make(chan Stream)
	go func(ch chan<- Stream) {
		// tells to every channel listener to stop listening
		defer close(ch)

		for {
			conn, err := lst.Accept()
			if err != nil {
				fmt.Printf("[ERROR] while accepting incomming connection: %s", err.Error())
				return
			}

			ch <- Stream{tcpConn: conn, remote: conn.RemoteAddr()}
		}
	}(connCh)

	return connCh, nil
}

func Dial(peerAddrPort string) (*Stream, error) {
	addrPort, err := netip.ParseAddrPort(peerAddrPort)
	if err != nil {
		return nil, fmt.Errorf("parsing addr and port: %w", err)
	}

	addr := net.TCPAddrFromAddrPort(addrPort)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("whiel dialing: %w", err)
	}

	return &Stream{
		tcpConn: conn,
	}, nil
}

func (s *Stream) Send(buff []byte) error {
	toBeSent := len(buff)
	sent := 0

	for sent != toBeSent {
		n, err := s.tcpConn.Write(buff)
		if err != nil {
			return fmt.Errorf("sent %d bytes, error while writing: %w", n, err)
		}
		sent += n
		fmt.Printf("sent %d bytes (total %d)...\n", sent, toBeSent)
	}

	return nil
}

func (s *Stream) WaitResponse(buff []byte) (int, error) {
	return s.tcpConn.Read(buff)
}
