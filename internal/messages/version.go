package messages

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/netip"

	"github.com/EclesioMeloJunior/btc-handshake/internal/codec"
)

var ErrUnexpectedEncodedRelay = errors.New("unexpected encoded relay")

type VersionOpt func(*Version)

func WithNumber(num uint32) VersionOpt {
	return func(v *Version) {
		v.Number = num
	}
}

func WithServices(services uint64) VersionOpt {
	return func(v *Version) {
		v.Services = services
	}
}

func WithTimestamp(ts int64) VersionOpt {
	return func(v *Version) {
		v.Timestamp = ts
	}
}

func WithAddrFrom(addr string, port uint16, services uint64) VersionOpt {
	return func(v *Version) {
		v.AddrFrom = NetworkAddress{
			Services: services,
			IpV6V4:   netip.MustParseAddr(addr),
			Port:     port,
		}
	}
}

func WithAddrRecv(addr string, port uint16, services uint64) VersionOpt {
	return func(v *Version) {
		v.AddrRecv = NetworkAddress{
			Services: services,
			IpV6V4:   netip.MustParseAddr(addr),
			Port:     port,
		}
	}
}

func WithNonce(nonce uint64) VersionOpt {
	return func(v *Version) {
		v.Nonce = nonce
	}
}

func WithUserAgent(userAgent string) VersionOpt {
	return func(v *Version) {
		v.UserAgent = userAgent
	}
}

func WithStartHeight(height uint32) VersionOpt {
	return func(v *Version) {
		v.StartHeight = height
	}
}

func AsRelay() VersionOpt {
	return func(v *Version) {
		v.Relay = true
	}
}

type Version struct {
	Number      uint32
	Services    uint64
	Timestamp   int64
	AddrRecv    NetworkAddress
	AddrFrom    NetworkAddress
	Nonce       uint64
	UserAgent   string
	StartHeight uint32
	Relay       bool
}

func NewVersion(opts ...VersionOpt) *Version {
	version := new(Version)
	for _, opt := range opts {
		opt(version)
	}

	return version
}

func (v *Version) Encode() ([]byte, error) {
	encodedVersion := make([]byte, 4)
	binary.LittleEndian.PutUint32(encodedVersion, v.Number)

	encodedServices := make([]byte, 8)
	binary.LittleEndian.PutUint64(encodedServices, v.Services)

	encodedTs := make([]byte, 8)
	binary.LittleEndian.PutUint64(encodedTs, uint64(v.Timestamp))

	encodedRcv, err := v.AddrRecv.Encode()
	if err != nil {
		return nil, fmt.Errorf("while encoding receiver addr: %w", err)
	}

	encodedFrom, err := v.AddrFrom.Encode()
	if err != nil {
		return nil, fmt.Errorf("while encoding from addr: %w", err)
	}

	encodedNonce := make([]byte, 8)
	binary.LittleEndian.PutUint64(encodedNonce, v.Nonce)

	encodedUA := codec.EncodeString(v.UserAgent)
	if len(encodedUA) == 0 {
		encodedUA = []byte{0x00}
	}

	encodedStartHeight := make([]byte, 4)
	binary.LittleEndian.PutUint32(encodedStartHeight, v.StartHeight)

	encodedRelay := []byte{0x00}
	if v.Relay {
		encodedRelay = []byte{0x01}
	}

	return bytes.Join([][]byte{
		encodedVersion,
		encodedServices,
		encodedTs,
		encodedRcv,
		encodedFrom,
		encodedNonce,
		encodedUA,
		encodedStartHeight,
		encodedRelay,
	}, nil), nil
}

func (v *Version) Decode(r io.Reader) error {
	enc := make([]byte, 4)
	_, err := r.Read(enc)
	if err != nil {
		return fmt.Errorf("while reading version number: %w", err)
	}
	v.Number = binary.LittleEndian.Uint32(enc)

	enc = make([]byte, 8)
	_, err = r.Read(enc)
	if err != nil {
		return fmt.Errorf("while reading services: %w", err)
	}
	v.Services = binary.LittleEndian.Uint64(enc)

	_, err = r.Read(enc)
	if err != nil {
		return fmt.Errorf("while reading timestamp: %w", err)
	}
	v.Timestamp = int64(binary.LittleEndian.Uint64(enc))

	err = v.AddrRecv.Decode(r)
	if err != nil {
		return fmt.Errorf("while decoding receiver addr: %w", err)
	}

	err = v.AddrFrom.Decode(r)
	if err != nil {
		return fmt.Errorf("while decoding from addr: %w", err)
	}

	_, err = r.Read(enc)
	if err != nil {
		return fmt.Errorf("while reading nonce: %w", err)
	}
	v.Nonce = binary.LittleEndian.Uint64(enc)

	v.UserAgent, err = codec.DecodeVarString(r)
	if err != nil {
		return fmt.Errorf("while decoding user agent: %w", err)
	}

	enc = make([]byte, 4)
	_, err = r.Read(enc)
	if err != nil {
		return fmt.Errorf("while reading start height: %w", err)
	}
	v.StartHeight = binary.LittleEndian.Uint32(enc)

	enc = make([]byte, 1)
	_, err = r.Read(enc)
	if err != nil {
		return fmt.Errorf("while reading realy: %w", err)
	}

	switch enc[0] {
	case 0x00:
		v.Relay = false
	case 0x01:
		v.Relay = true
	default:
		return fmt.Errorf("%w, byte must be 0x00 (false) or 0x01 (true), got: %v", ErrUnexpectedEncodedRelay, enc[0])
	}

	return nil
}
