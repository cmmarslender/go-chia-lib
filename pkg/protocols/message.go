package protocols

import (
	"github.com/cmmarslender/go-chia-lib/pkg/streamable"
)

// Message is a protocol message
type Message struct {
	ProtocolMessageType ProtocolMessageType `streamable:""`
	ID                  *uint16             `streamable:"optional"`
	Data                []byte              `streamable:""`
}

// MakeMessage makes a new Message with the given data
func MakeMessage(messageType ProtocolMessageType, data interface{}) (*Message, error) {
	msg := &Message{
		ProtocolMessageType: messageType,
	}

	dataBytes, err := streamable.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg.Data = dataBytes

	return msg, nil
}

// MakeMessageBytes calls MakeMessage and converts everything down to bytes
func MakeMessageBytes(messageType ProtocolMessageType, data interface{}) ([]byte, error) {
	msg, err := MakeMessage(messageType, data)
	if err != nil {
		return nil, err
	}

	return streamable.Marshal(msg)
}

// DecodeMessage is a helper function to quickly decode bytes to Message
func DecodeMessage(bytes []byte) (*Message, error) {
	msg := &Message{}

	err := streamable.Unmarshal(bytes, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// DecodeMessageData decodes a message.data into the given interface
func DecodeMessageData(bytes []byte, v interface{}) error {
	msg, err := DecodeMessage(bytes)
	if err != nil {
		return err
	}

	return streamable.Unmarshal(msg.Data, v)
}
