package protocols

import (
	"github.com/cmmarslender/go-chia-lib/pkg/types"
)

// RequestPeers is an empty struct
type RequestPeers struct{}

// RespondPeers is the format for the request_peers response
type RespondPeers struct {
	PeerList []types.TimestampedPeerInfo `streamable:""`
}
