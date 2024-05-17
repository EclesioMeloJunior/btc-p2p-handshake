package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var ErrWronglyEncodedVarint = errors.New("wrongly encoded varint")
var ErrUnexpectedReadSize = errors.New("unexpected read size")

type Encodeable interface {
	fmt.Stringer

	Encode() ([]byte, error)
	Decode(io.Reader) error
}

// LittleEndianPutUint32 similar to binary.LittleEndian this function
// takes a io.Writer instead of a []byte as argument and returns an error
// if could not write the encoded little endian number
func LittleEndianPutUint32(r io.Writer, v uint32) error {
	_, err := r.Write([]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)})
	return err
}

func EncodeToVarint(val uint64) []byte {
	switch {
	case val < 0xFD:
		return []byte{uint8(val)}
	case val <= 0xFF_FF:
		enc := make([]byte, 3)
		enc[0] = 0xFD
		binary.LittleEndian.PutUint16(enc[1:], uint16(val))
		return enc
	case val <= 0xFFFF_FFFF:
		enc := make([]byte, 5)
		enc[0] = 0xFE
		binary.LittleEndian.PutUint32(enc[1:], uint32(val))
		return enc
	default:
		enc := make([]byte, 9)
		enc[0] = 0xFF
		binary.LittleEndian.PutUint64(enc[1:], val)
		return enc
	}
}

func DecodeFromVarint(r io.Reader) (uint64, error) {
	fst := make([]byte, 1)
	_, err := r.Read(fst)
	if err != nil {
		return 0, err
	}

	if fst[0] < 0xFD {
		return uint64(fst[0]), nil
	}

	var rest []byte
	switch fst[0] {
	case 0xFD:
		rest = make([]byte, 2)
		_, err = r.Read(rest)
		if err != nil {
			return 0, err
		}

		return uint64(binary.LittleEndian.Uint16(rest)), nil
	case 0xFE:
		rest = make([]byte, 4)
		_, err = r.Read(rest)
		if err != nil {
			return 0, err
		}

		return uint64(binary.LittleEndian.Uint32(rest)), nil
	case 0xFF:
		rest = make([]byte, 8)
		_, err = r.Read(rest)
		if err != nil {
			return 0, err
		}

		return binary.LittleEndian.Uint64(rest), nil
	default:
		return 0, fmt.Errorf("%w: unsupported pre-appended byte %d", ErrWronglyEncodedVarint, fst[0])
	}
}

func EncodeString(s string) []byte {
	if len(s) == 0 {
		return nil
	}

	encLenght := EncodeToVarint(uint64(len(s)))
	return bytes.Join([][]byte{encLenght, []byte(s)}, nil)
}

func DecodeVarString(r io.Reader) (string, error) {
	strLen, err := DecodeFromVarint(r)
	if err != nil {
		return "", fmt.Errorf("decoding string length: %w", err)
	}

	full := make([]byte, strLen)
	n, err := r.Read(full)
	if err != nil {
		return "", fmt.Errorf("while reading encoded string: %w", err)
	}

	if uint64(n) != strLen {
		return "", fmt.Errorf("%w, expected: %d, actual: %d", ErrUnexpectedReadSize, strLen, n)
	}

	return string(full), nil
}
