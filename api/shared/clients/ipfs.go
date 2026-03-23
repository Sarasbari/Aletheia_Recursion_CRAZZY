package clients

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"

	shell "github.com/ipfs/go-ipfs-api"
)

type IPFSClient struct {
	shell *shell.Shell
}

func NewIPFSClient(endpoint string) *IPFSClient {
	if endpoint == "" {
		return &IPFSClient{}
	}
	return &IPFSClient{shell: shell.NewShell(endpoint)}
}

func (c *IPFSClient) Upload(_ context.Context, payload []byte) (string, error) {
	if c.shell == nil {
		sum := sha256.Sum256(payload)
		return "local-" + hex.EncodeToString(sum[:])[:32], nil
	}
	return c.shell.Add(bytes.NewReader(payload))
}
