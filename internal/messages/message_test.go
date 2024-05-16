package messages_test

import (
	"bytes"
	"fmt"
	"io"
	"net/netip"
	"testing"

	"github.com/EclesioMeloJunior/btc-handshake/internal/messages"
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

	m := &messages.Message{
		Magic:   messages.MagicMain,
		Command: []byte("version"),
		Payload: &payload,
	}

	enc, err := m.Encode()
	require.NoError(t, err)
	fmt.Printf("0x%x\n", enc)
}

func TestNetworkAddressEncoding(t *testing.T) {
	// taken from an example at: https://en.bitcoin.it/wiki/Protocol_documentation#Network_address
	testEncoded := []byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // NODE_NETWORK service
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x0A, 0x00, 0x00, 0x01, // ip addr
		0x20, 0x8D,
	}

	netaddr := &messages.NetworkAddress{}
	err := netaddr.Decode(bytes.NewReader(testEncoded))
	require.NoError(t, err)

	expectedNetAddr := &messages.NetworkAddress{
		Services: 1,
		IpV6V4:   netip.MustParseAddr("10.0.0.1"),
		Port:     8333,
	}
	require.Equal(t, expectedNetAddr, expectedNetAddr)

	encoded, err := expectedNetAddr.Encode()
	require.NoError(t, err)
	require.Equal(t, testEncoded, encoded)
}
