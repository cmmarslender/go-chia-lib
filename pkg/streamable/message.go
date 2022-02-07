package streamable

// ProtocolMessageType corresponds to ProtocolMessageTypes in Chia
type ProtocolMessageType uint8

const (
	// ProtocolMessageTypeHandshake Handshake
	ProtocolMessageTypeHandshake ProtocolMessageType = 1

	// there are many more of these in Chia - only listing the ones current is use for now

	// ProtocolMessageTypeRequestPeers request_peers
	ProtocolMessageTypeRequestPeers ProtocolMessageType = 43

	// ProtocolMessageTypeRespondPeers respond_peers
	ProtocolMessageTypeRespondPeers ProtocolMessageType = 44
)

// Message is a protocol message
type Message struct {
	ProtocolMessageType ProtocolMessageType `streamable:""`
	ID                  *uint16             `streamable:"optional"`
	Data                []byte              `streamable:""`
}
