package messages

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/netip"

	"github.com/EclesioMeloJunior/btc-handshake/internal/codec"
)

var (
	ErrCommandBytesOverflow  = errors.New("message command bytes overflow")
	ErrFailedToEncodePayload = errors.New("failed to encode payload")
	ErrFailedToEncodeMessage = errors.New("failed to encode message")
)

var IPV6Default = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF}

type NodeService = uint64

const (
	NodeNetwork        NodeService = 1
	NodeGetUtxo        NodeService = 2
	NodeBloom          NodeService = 4
	NodeWitness        NodeService = 8
	NodeXThin          NodeService = 16
	NodeCompactFilters NodeService = 64
	NodeNetworkLimited NodeService = 1024
)

type Magic uint32

const (
	MagicMain           Magic = 0xD9B4BEF9
	MagicTestNetRegTest Magic = 0xDAB5BFFA
	MagicTestNet3       Magic = 0x0709110B
	MagicSignet         Magic = 0x40CF030A
	MagicNameCoin       Magic = 0xFEB4BEF9
)

type Service uint64

const MessageCapWithoutPayloadInBytes = 24

var _ codec.Encodeable = (*Message)(nil)
var _ codec.Encodeable = (*NetworkAddress)(nil)
var _ codec.Encodeable = (*Version)(nil)

type EmptyPayload struct{}

func (EmptyPayload) Encode() ([]byte, error) {
	return nil, nil
}

func (EmptyPayload) Decode(_ io.Reader) error {
	return nil
}

type Message struct {
	Magic   Magic
	Command []byte
	Payload codec.Encodeable
}

func NewMainMessage(command []byte, payload codec.Encodeable) Message {
	return Message{
		Magic:   MagicMain,
		Command: command,
		Payload: payload,
	}
}

func (m Message) Encode() ([]byte, error) {
	encodedPayload, err := m.Payload.Encode()
	if err != nil {
		return nil, errors.Join(ErrFailedToEncodePayload, err)
	}

	// initialize a buffer with the exact message capacity + payload size
	// check: https://en.bitcoin.it/wiki/Protocol_documentation#Message_structure
	encodedMessageOutput := bytes.NewBuffer(make([]byte, 0, MessageCapWithoutPayloadInBytes+len(encodedPayload)))

	if err := codec.LittleEndianPutUint32(encodedMessageOutput, uint32(m.Magic)); err != nil {
		return nil, errors.Join(ErrFailedToEncodeMessage, err)
	}

	n, err := encodedMessageOutput.Write(m.Command)
	if err != nil {
		return nil, errors.Join(ErrFailedToEncodeMessage, err)
	}

	if n > 12 {
		return nil, fmt.Errorf("%w: 12 is the current limit", ErrCommandBytesOverflow)
	} else if n < 12 {
		// padding with NULL
		encodedMessageOutput.Write(make([]byte, 12-n))
	}

	if err := codec.LittleEndianPutUint32(encodedMessageOutput, uint32(len(encodedPayload))); err != nil {
		return nil, errors.Join(ErrFailedToEncodeMessage, err)
	}

	if _, err := encodedMessageOutput.Write(checksum(encodedPayload)); err != nil {
		return nil, errors.Join(ErrFailedToEncodeMessage, err)
	}

	if _, err := encodedMessageOutput.Write(encodedPayload); err != nil {
		return nil, errors.Join(ErrFailedToEncodeMessage, err)
	}

	return encodedMessageOutput.Bytes(), nil
}

func checksum(value []byte) []byte {
	fst := sha256.Sum256(value)
	snd := sha256.Sum256(fst[:])
	return snd[:4]
}

func (m *Message) Decode(_ io.Reader) error {
	return nil
}

type NetworkAddress struct {
	// the protocol specifies about a field called
	// `time` which is not present when the net_addr
	// is part of a version message.
	// However my big question is: the protocol specify
	// this field as uint32_i however the time (as timestamp)
	// is a int64_i, so how is possible to encode w/o data loss?
	//
	// time uint32

	Services uint64
	IpV6V4   netip.Addr
	Port     uint16
}

func (n *NetworkAddress) Encode() ([]byte, error) {
	enc := make([]byte, 26)
	binary.LittleEndian.PutUint64(enc[:8], n.Services)

	// since we're using ipv4 here, we will set the first
	// 12 bytes (from 8 to 20) and the next 4 bytes are the
	// actual encoded IPV4
	copy(enc[8:20], IPV6Default)
	copy(enc[20:24], n.IpV6V4.AsSlice())
	binary.BigEndian.PutUint16(enc[24:], n.Port)
	return enc, nil
}

func (n *NetworkAddress) Decode(r io.Reader) error {
	enc := make([]byte, 8)
	_, err := r.Read(enc)
	if err != nil {
		return fmt.Errorf("reading services: %w", err)
	}
	n.Services = binary.LittleEndian.Uint64(enc)

	enc = make([]byte, 16)
	_, err = r.Read(enc)
	if err != nil {
		return fmt.Errorf("reading ipv6/v4 address: %w", err)
	}

	err = n.IpV6V4.UnmarshalBinary(enc[12:])
	if err != nil {
		return fmt.Errorf("unmarshaling ipv6/v4")
	}

	enc = make([]byte, 2)
	_, err = r.Read(enc)
	if err != nil {
		return fmt.Errorf("reading port: %w", err)
	}

	n.Port = binary.BigEndian.Uint16(enc)
	return nil
}
