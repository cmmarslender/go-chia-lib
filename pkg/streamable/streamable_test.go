package streamable_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cmmarslender/go-chia-lib/pkg/streamable"
)

const (
	// This is the encoded version of the following:
	//test_message = make_msg(
	//	ProtocolMessageTypes.handshake,
	//	"This is a sample message to decode".encode(encoding = 'UTF-8', errors = 'string')
	//)
	encodedHex string = "0100000000225468697320697320612073616d706c65206d65737361676520746f206465636f6465"
)

func TestUnmarshal(t *testing.T) {
	// Hex to bytes
	encodedBytes, err := hex.DecodeString(encodedHex)
	assert.NoError(t, err)

	// test that nil is not accepted
	err = streamable.Unmarshal(encodedBytes, nil)
	assert.Error(t, err)

	msg := &streamable.Message{
		ID:                  nil,
		ProtocolMessageType: 0,
		Data:                nil,
	}

	// Test that pointers are required
	err = streamable.Unmarshal(encodedBytes, *msg)
	assert.Error(t, err)

	err = streamable.Unmarshal(encodedBytes, msg)

	assert.NoError(t, err)
	assert.Equal(t, streamable.ProtocolMessageTypeHandshake, msg.ProtocolMessageType)
	assert.Nil(t, msg.ID) // @TODO need a testcase that also parses this field
	assert.Equal(t, []byte("This is a sample message to decode"), msg.Data)
}

func TestMarshal(t *testing.T) {
	msg := &streamable.Message{
		ProtocolMessageType: streamable.ProtocolMessageTypeHandshake,
		ID:                  nil,
		Data:                []byte("This is a sample message to decode"),
	}

	streamable.Marshal(msg)
}
