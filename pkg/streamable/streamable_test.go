package streamable_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cmmarslender/go-chia-lib/pkg/streamable"
	"github.com/cmmarslender/go-chia-lib/pkg/util"
)

const (
	//Message(
	//	uint8(ProtocolMessageTypes.handshake.value),
	//	None,
	//	bytes("This is a sample message to decode".encode(encoding = 'UTF-8', errors = 'string'))
	//)
	encodedHex1 string = "0100000000225468697320697320612073616d706c65206d65737361676520746f206465636f6465"

	//Message(
	//	uint8(ProtocolMessageTypes.handshake.value),
	//	uint16(35256),
	//	bytes("This is a sample message to decode".encode(encoding = 'UTF-8', errors = 'string'))
	//)
	encodedHex2 string = "010189b8000000225468697320697320612073616d706c65206d65737361676520746f206465636f6465"
)

func TestUnmarshal_Message1(t *testing.T) {
	// Hex to bytes
	encodedBytes, err := hex.DecodeString(encodedHex1)
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
	assert.Nil(t, msg.ID)
	assert.Equal(t, []byte("This is a sample message to decode"), msg.Data)
}

func TestMarshal_Message1(t *testing.T) {
	encodedBytes, err := hex.DecodeString(encodedHex1)
	assert.NoError(t, err)

	msg := &streamable.Message{
		ProtocolMessageType: streamable.ProtocolMessageTypeHandshake,
		ID:                  nil,
		Data:                []byte("This is a sample message to decode"),
	}

	bytes, err := streamable.Marshal(msg)

	assert.NoError(t, err)
	assert.Equal(t, encodedBytes, bytes)
}

func TestUnmarshal_Message2(t *testing.T) {
	// Hex to bytes
	encodedBytes, err := hex.DecodeString(encodedHex2)
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
	assert.Equal(t, util.PtrUint16(35256), msg.ID)
	assert.Equal(t, []byte("This is a sample message to decode"), msg.Data)
}

func TestMarshal_Message2(t *testing.T) {
	encodedBytes, err := hex.DecodeString(encodedHex2)
	assert.NoError(t, err)

	msg := &streamable.Message{
		ProtocolMessageType: streamable.ProtocolMessageTypeHandshake,
		ID:                  util.PtrUint16(35256),
		Data:                []byte("This is a sample message to decode"),
	}

	bytes, err := streamable.Marshal(msg)

	assert.NoError(t, err)
	assert.Equal(t, encodedBytes, bytes)
}
