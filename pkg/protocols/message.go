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
