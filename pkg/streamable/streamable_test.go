package streamable_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cmmarslender/go-chia-lib/pkg/protocols"
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

	//Message(
	//	uint8(ProtocolMessageTypes.handshake.value),
	//	None,
	//	Handshake(
	//      "mainnet",
	//      "0.0.33",
	//      "1.2.11",
	//      uint16(8444),
	//      uint8(1),
	//      [(uint16(Capability.BASE.value), "1")],
	//  )
	//)
	encodedHexHandshake string = "01000000002d000000076d61696e6e657400000006302e302e333300000006312e322e313120fc010000000100010000000131"
)

func TestUnmarshal_Message1(t *testing.T) {
	// Hex to bytes
	encodedBytes, err := hex.DecodeString(encodedHex1)
	assert.NoError(t, err)

	// test that nil is not accepted
	err = streamable.Unmarshal(encodedBytes, nil)
	assert.Error(t, err)

	msg := &protocols.Message{
		ID:                  nil,
		ProtocolMessageType: 0,
		Data:                nil,
	}

	// Test that pointers are required
	err = streamable.Unmarshal(encodedBytes, *msg)
	assert.Error(t, err)

	err = streamable.Unmarshal(encodedBytes, msg)

	assert.NoError(t, err)
	assert.Equal(t, protocols.ProtocolMessageTypeHandshake, msg.ProtocolMessageType)
	assert.Nil(t, msg.ID)
	assert.Equal(t, []byte("This is a sample message to decode"), msg.Data)
}

func TestMarshal_Message1(t *testing.T) {
	encodedBytes, err := hex.DecodeString(encodedHex1)
	assert.NoError(t, err)

	msg := &protocols.Message{
		ProtocolMessageType: protocols.ProtocolMessageTypeHandshake,
		ID:                  nil,
		Data:                []byte("This is a sample message to decode"),
	}

	bytes, err := streamable.Marshal(msg)

	assert.NoError(t, err)
	assert.Equal(t, encodedBytes, bytes)
}

// Unmarshals fully then remarshals to ensure we can go back and forth
func TestUnmarshal_Remarshal_Message1(t *testing.T) {
	encodedBytes, err := hex.DecodeString(encodedHex1)
	assert.NoError(t, err)

	msg := &protocols.Message{}

	err = streamable.Unmarshal(encodedBytes, msg)
	assert.NoError(t, err)

	// Remarshal and check against original bytes
	reencodedBytes, err := streamable.Marshal(msg)
	assert.NoError(t, err)
	assert.Equal(t, encodedBytes, reencodedBytes)
}

func TestUnmarshal_Message2(t *testing.T) {
	// Hex to bytes
	encodedBytes, err := hex.DecodeString(encodedHex2)
	assert.NoError(t, err)

	// test that nil is not accepted
	err = streamable.Unmarshal(encodedBytes, nil)
	assert.Error(t, err)

	msg := &protocols.Message{
		ID:                  nil,
		ProtocolMessageType: 0,
		Data:                nil,
	}

	// Test that pointers are required
	err = streamable.Unmarshal(encodedBytes, *msg)
	assert.Error(t, err)

	err = streamable.Unmarshal(encodedBytes, msg)

	assert.NoError(t, err)
	assert.Equal(t, protocols.ProtocolMessageTypeHandshake, msg.ProtocolMessageType)
	assert.Equal(t, util.PtrUint16(35256), msg.ID)
	assert.Equal(t, []byte("This is a sample message to decode"), msg.Data)
}

func TestMarshal_Message2(t *testing.T) {
	encodedBytes, err := hex.DecodeString(encodedHex2)
	assert.NoError(t, err)

	msg := &protocols.Message{
		ProtocolMessageType: protocols.ProtocolMessageTypeHandshake,
		ID:                  util.PtrUint16(35256),
		Data:                []byte("This is a sample message to decode"),
	}

	bytes, err := streamable.Marshal(msg)

	assert.NoError(t, err)
	assert.Equal(t, encodedBytes, bytes)
}

// Unmarshals fully then remarshals to ensure we can go back and forth
func TestUnmarshal_Remarshal_Message2(t *testing.T) {
	encodedBytes, err := hex.DecodeString(encodedHex2)
	assert.NoError(t, err)

	msg := &protocols.Message{}

	err = streamable.Unmarshal(encodedBytes, msg)
	assert.NoError(t, err)

	// Remarshal and check against original bytes
	reencodedBytes, err := streamable.Marshal(msg)
	assert.NoError(t, err)
	assert.Equal(t, encodedBytes, reencodedBytes)
}

func TestUnmarshal_Handshake(t *testing.T) {
	// Hex to bytes
	encodedBytes, err := hex.DecodeString(encodedHexHandshake)
	assert.NoError(t, err)

	msg := &protocols.Message{}

	err = streamable.Unmarshal(encodedBytes, msg)

	assert.NoError(t, err)
	assert.Equal(t, protocols.ProtocolMessageTypeHandshake, msg.ProtocolMessageType)
	assert.Nil(t, msg.ID)

	// No decode the handshake portion
	handshake := &protocols.Handshake{}

	//	Handshake(
	//      "mainnet",
	//      "0.0.33",
	//      "1.2.11",
	//      uint16(8444),
	//      uint8(1),
	//      [(uint16(Capability.BASE.value), "1")],
	//  )

	err = streamable.Unmarshal(msg.Data, handshake)
	assert.NoError(t, err)
	assert.Equal(t, "mainnet", handshake.NetworkID)
	assert.Equal(t, "0.0.33", handshake.ProtocolVersion)
	assert.Equal(t, "1.2.11", handshake.SoftwareVersion)
	assert.Equal(t, uint16(8444), handshake.ServerPort)
	assert.Equal(t, protocols.NodeTypeFullNode, handshake.NodeType)
	assert.IsType(t, []protocols.Capability{}, handshake.Capabilities)
	assert.Len(t, handshake.Capabilities, 1)

	// Test each capability item
	cap1 := handshake.Capabilities[0]

	assert.Equal(t, protocols.CapabilityTypeBase, cap1.Capability)
	assert.Equal(t, "1", cap1.Value)
}

// Unmarshals fully then remarshals to ensure we can go back and forth
func TestUnmarshal_Remarshal_Handshake(t *testing.T) {
	encodedBytes, err := hex.DecodeString(encodedHexHandshake)
	assert.NoError(t, err)

	msg := &protocols.Message{}

	err = streamable.Unmarshal(encodedBytes, msg)
	assert.NoError(t, err)

	handshake := &protocols.Handshake{}

	err = streamable.Unmarshal(msg.Data, handshake)
	assert.NoError(t, err)

	// Remarshal and check against original bytes
	reencodedBytes, err := protocols.MakeMessageBytes(msg.ProtocolMessageType, handshake)
	assert.NoError(t, err)
	assert.Equal(t, encodedBytes, reencodedBytes)
}
