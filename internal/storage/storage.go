package storage

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sh "github.com/ipfs/go-ipfs-api"
)

// ImageStorage defines the interface for image storage.
type ImageStorage interface {
	UploadImage(ctx context.Context, data []byte) (string, error)
}

// MetadataStorage defines the interface for metadata storage.
type MetadataStorage interface {
	StoreMetadata(ctx context.Context, metadata interface{}) error
}

// IPFSStorage implements Storage for IPFS.
type IPFSStorage struct {
	shell *sh.Shell
}

func NewIPFSStorage(url string) *IPFSStorage {
	return &IPFSStorage{shell: sh.NewShell(url)}
}

func (s *IPFSStorage) UploadImage(ctx context.Context, data []byte) (string, error) {
	cid, err := s.shell.Add(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to upload to IPFS: %w", err)
	}
	return cid, nil
}

// BigchainDBStorage implements MetadataStorage for BigchainDB.
type BigchainDBStorage struct {
	nodeURL    string
	privateKey ed25519.PrivateKey // Ed25519 private key
	publicKey  ed25519.PublicKey  // Ed25519 public key
}

func NewBigchainDBStorage(nodeURL string, seedHex string) (*BigchainDBStorage, error) {
	seed, err := hex.DecodeString(seedHex)
	if err != nil {
		return nil, fmt.Errorf("invalid BigchainDB seed: %w", err)
	}
	if len(seed) != 32 {
		return nil, fmt.Errorf("BigchainDB seed must be 32 bytes")
	}

	priv := ed25519.NewKeyFromSeed(seed)
	pub := priv.Public().(ed25519.PublicKey)

	return &BigchainDBStorage{
		nodeURL:    nodeURL,
		privateKey: priv,
		publicKey:  pub,
	}, nil
}

func (s *BigchainDBStorage) StoreMetadata(ctx context.Context, metadata interface{}) error {
	// Construct a BigchainDB CREATE transaction
	tx := map[string]interface{}{
		"operation": "CREATE",
		"version":   "2.0",
		"asset": map[string]interface{}{
			"data": metadata,
		},
		"metadata": nil,
		"outputs": []interface{}{
			map[string]interface{}{
				"amount": "1",
				"condition": map[string]interface{}{
					"details": map[string]interface{}{
						"public_key": hex.EncodeToString(s.publicKey),
						"type":       "ed25519-sha-256",
					},
					"uri": fmt.Sprintf("ni:///sha-256;%s?fpt=ed25519-sha-256&cost=131072", hex.EncodeToString(s.publicKey)),
				},
				"public_keys": []string{hex.EncodeToString(s.publicKey)},
			},
		},
		"inputs": []interface{}{
			map[string]interface{}{
				"fulfillment":   nil,
				"fulfills":      nil,
				"owners_before": []string{hex.EncodeToString(s.publicKey)},
			},
		},
	}

	// 1. Serialize and hash the transaction to get the ID (excluding some fields usually)
	// For simplicity, we use the serialized JSON of the entire structure
	txJSON, _ := json.Marshal(tx)
	hash := sha256.Sum256(txJSON)
	txID := hex.EncodeToString(hash[:])
	tx["id"] = txID

	// 2. Sign the transaction (simplified fulfillment)
	signature := ed25519.Sign(s.privateKey, hash[:])
	fulfillment := hex.EncodeToString(signature)
	tx["inputs"].([]interface{})[0].(map[string]interface{})["fulfillment"] = fulfillment

	// 3. Post to BigchainDB
	finalJSON, _ := json.Marshal(tx)
	resp, err := http.Post(s.nodeURL+"/api/v1/transactions", "application/json", bytes.NewBuffer(finalJSON))
	if err != nil {
		return fmt.Errorf("failed to post to BigchainDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("BigchainDB error: %s (status: %d)", string(body), resp.StatusCode)
	}

	return nil
}
