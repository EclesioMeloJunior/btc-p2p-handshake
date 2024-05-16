package codec_test

import (
	"bytes"
	"testing"

	"github.com/EclesioMeloJunior/btc-handshake/internal/codec"
	"github.com/stretchr/testify/require"
)

func TestVarint(t *testing.T) {
	cases := []struct {
		value   uint64
		encoded []byte
	}{
		{
			value:   10,
			encoded: []byte{0x0a},
		},
		{
			value:   61951,
			encoded: []byte{0xFD, 0xFF, 0xF1},
		},
		{
			value:   61951,
			encoded: []byte{0xFD, 0xFF, 0xF1},
		},
		{
			value:   4289466879,
			encoded: []byte{0xFE, 0xFF, 0x11, 0xAC, 0xFF},
		},
		{
			value:   ^uint64(0),
			encoded: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
	}

	for _, tt := range cases {
		enc := codec.EncodeToVarint(tt.value)
		require.Equal(t, tt.encoded, enc)

		dec, err := codec.DecodeFromVarint(bytes.NewBuffer(enc))
		require.NoError(t, err)
		require.Equal(t, tt.value, dec)
	}
}
