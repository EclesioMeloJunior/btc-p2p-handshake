package network_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/EclesioMeloJunior/btc-handshake/internal/network"
	"github.com/stretchr/testify/require"
)

type DummyEncodable []byte

func (d *DummyEncodable) Encode() ([]byte, error) {
	return *d, nil
}

func (*DummyEncodable) Decode(_ io.Reader) error {
	return nil
}

func TestMessageEncoding(t *testing.T) {
	payload := DummyEncodable(make([]byte, 100))

	m := &network.Message{
		Magic:   network.Main,
		Command: []byte("version"),
		Payload: &payload,
	}

	enc, err := m.Encode()
	require.NoError(t, err)
	fmt.Printf("0x%x\n", enc)
}
