package streamable

// NodeType is the type of peer (farmer, full node, etc)
// @todo This is also duplicated in the RPC lib. Probably a good candidate for the go-lib
type NodeType uint8

const (
	// NodeTypeFullNode NodeType for full node
	NodeTypeFullNode NodeType = 1

	// NodeTypeHarvester NodeType for Harvester
	NodeTypeHarvester NodeType = 2

	// NodeTypeFarmer NodeType for Farmer
	NodeTypeFarmer NodeType = 3

	// NodeTypeTimelord NodeType for Timelord
	NodeTypeTimelord NodeType = 4

	// NodeTypeIntroducer NodeType for Introducer
	NodeTypeIntroducer NodeType = 5

	// NodeTypeWallet NodeType for Wallet
	NodeTypeWallet NodeType = 6
)

// Handshake is a handshake message
type Handshake struct {
	NetworkID       string   `streamable:""`
	ProtocolVersion string   `streamable:""`
	SoftwareVersion string   `streamable:""`
	ServerPort      uint16   `streamable:""`
	NodeType        NodeType `streamable:""`
	// Capabilities List[Tuple[uint16, str]] // @TODO Support this
}
