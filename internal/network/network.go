package network

import (
	"fmt"
	"net"
	"net/netip"
)

type Server struct {
	Err      chan error
	peerAddr *net.TCPAddr
	tcpConn  *net.TCPConn
}

func New(peerAddrPort string) (*Server, error) {
	addrPort, err := netip.ParseAddrPort(peerAddrPort)
	if err != nil {
		return nil, fmt.Errorf("parsing addr and port: %w", err)
	}

	addr := net.TCPAddrFromAddrPort(addrPort)

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("whiel dialing: %w", err)
	}

	return &Server{
		peerAddr: addr,
		tcpConn:  conn,
	}, nil
}

func (s *Server) Send(buff []byte) error {
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

func (s *Server) WaitResponse(buff []byte) (int, error) {
	return s.tcpConn.Read(buff)
}
