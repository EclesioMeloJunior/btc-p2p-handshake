package network

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

var (
	ErrCommandBytesOverflow  = errors.New("message command bytes overflow")
	ErrFailedToEncodePayload = errors.New("failed to encode payload")
	ErrFailedToEncodeMessage = errors.New("failed to encode message")
)

type Magic uint32

const (
	Main           Magic = 0xD9B4BEF9
	TestNetRegTest Magic = 0xDAB5BFFA
	TestNet3       Magic = 0x0709110B
	Signet         Magic = 0x40CF030A
	NameCoin       Magic = 0xFEB4BEF9
)

const MessageCapWithoutPayloadInBytes = 24

var _ Encodeable = (*Message)(nil)

type Message struct {
	Magic   Magic
	Command []byte
	Payload Encodeable
}

func (m *Message) Encode() ([]byte, error) {
	encodedPayload, err := m.Payload.Encode()
	if err != nil {
		return nil, errors.Join(ErrFailedToEncodePayload, err)
	}

	// initialize a buffer with the exact message capacity + payload size
	// check: https://en.bitcoin.it/wiki/Protocol_documentation#Message_structure
	encodedMessageOutput := bytes.NewBuffer(make([]byte, 0, MessageCapWithoutPayloadInBytes+len(encodedPayload)))

	if err := LittleEndianPutUint32(encodedMessageOutput, uint32(m.Magic)); err != nil {
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

	if err := LittleEndianPutUint32(encodedMessageOutput, uint32(len(encodedPayload))); err != nil {
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
