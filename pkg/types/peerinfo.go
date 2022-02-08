package types

// TimestampedPeerInfo contains information about peers with timestamps
type TimestampedPeerInfo struct {
	Host      string `streamable:""`
	Port      uint16 `streamable:""`
	timestamp uint64 `streamable:""`
}
