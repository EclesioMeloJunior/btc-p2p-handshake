package network

import "io"

type Encodeable interface {
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
