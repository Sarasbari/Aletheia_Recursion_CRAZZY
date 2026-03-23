package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"aletheia-api/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client interface {
	AnchorProof(ctx context.Context, record models.ProofRecord) (AnchorReceipt, error)
	AnchorVideoProof(ctx context.Context, record models.VideoRecord) (AnchorReceipt, error)
}

type AnchorReceipt struct {
	TxHash      string
	Timestamp   time.Time
	BlockNumber uint64
}

type PolygonClient struct {
	rpcURL          string
	privateKeyHex   string
	contractAddress string
	chainID         *big.Int
}

func NewPolygonClient(rpcURL, privateKeyHex, contractAddress, chainID string) *PolygonClient {
	cid := big.NewInt(137)
	if chainID != "" {
		if parsed, ok := new(big.Int).SetString(chainID, 10); ok {
			cid = parsed
		}
	}
	return &PolygonClient{
		rpcURL:          rpcURL,
		privateKeyHex:   privateKeyHex,
		contractAddress: contractAddress,
		chainID:         cid,
	}
}

func (p *PolygonClient) AnchorProof(ctx context.Context, record models.ProofRecord) (AnchorReceipt, error) {
	if p.rpcURL == "" || p.privateKeyHex == "" || p.contractAddress == "" {
		payload := fmt.Sprintf("%s:%s:%s:%s", record.SHA256, record.IPFSCID, record.MerkleRoot, record.Timestamp.Format(time.RFC3339Nano))
		h := crypto.Keccak256Hash([]byte(payload))
		now := time.Now().UTC()
		return AnchorReceipt{TxHash: "mock-" + h.Hex(), Timestamp: now, BlockNumber: uint64(now.Unix())}, nil
	}

	client, err := ethclient.DialContext(ctx, p.rpcURL)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("connect polygon rpc: %w", err)
	}
	defer client.Close()

	pk, err := crypto.HexToECDSA(p.privateKeyHex)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("decode private key: %w", err)
	}

	from := crypto.PubkeyToAddress(pk.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("suggest gas price: %w", err)
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("fetch latest block header: %w", err)
	}

	payload := map[string]string{
		"sha256":     record.SHA256,
		"ipfsCID":    record.IPFSCID,
		"merkleRoot": record.MerkleRoot,
		"timestamp":  record.Timestamp.Format(time.RFC3339Nano),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("marshal tx payload: %w", err)
	}

	to := common.HexToAddress(p.contractAddress)
	tx := types.NewTransaction(
		nonce,
		to,
		big.NewInt(0),
		300000,
		gasPrice,
		data,
	)

	signed, err := types.SignTx(tx, types.NewEIP155Signer(p.chainID), pk)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("sign tx: %w", err)
	}

	if err := client.SendTransaction(ctx, signed); err != nil {
		return AnchorReceipt{}, fmt.Errorf("send tx: %w", err)
	}

	return AnchorReceipt{
		TxHash:      signed.Hash().Hex(),
		Timestamp:   time.Unix(int64(header.Time), 0).UTC(),
		BlockNumber: header.Number.Uint64(),
	}, nil
}

func (p *PolygonClient) AnchorVideoProof(ctx context.Context, record models.VideoRecord) (AnchorReceipt, error) {
	if p.rpcURL == "" || p.privateKeyHex == "" || p.contractAddress == "" {
		payload := fmt.Sprintf("%s:%s:%s:%s:%s", record.VideoHash, record.AudioHash, record.VideoCID, record.AudioCID, record.Timestamp.Format(time.RFC3339Nano))
		h := crypto.Keccak256Hash([]byte(payload))
		now := time.Now().UTC()
		return AnchorReceipt{TxHash: "mock-" + h.Hex(), Timestamp: now, BlockNumber: uint64(now.Unix())}, nil
	}

	client, err := ethclient.DialContext(ctx, p.rpcURL)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("connect polygon rpc: %w", err)
	}
	defer client.Close()

	pk, err := crypto.HexToECDSA(p.privateKeyHex)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("decode private key: %w", err)
	}

	from := crypto.PubkeyToAddress(pk.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("suggest gas price: %w", err)
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("fetch latest block header: %w", err)
	}

	payload := map[string]string{
		"videoHash": record.VideoHash,
		"audioHash": record.AudioHash,
		"cidVideo":  record.VideoCID,
		"cidAudio":  record.AudioCID,
		"timestamp": record.Timestamp.Format(time.RFC3339Nano),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("marshal tx payload: %w", err)
	}

	to := common.HexToAddress(p.contractAddress)
	tx := types.NewTransaction(
		nonce,
		to,
		big.NewInt(0),
		300000,
		gasPrice,
		data,
	)

	signed, err := types.SignTx(tx, types.NewEIP155Signer(p.chainID), pk)
	if err != nil {
		return AnchorReceipt{}, fmt.Errorf("sign tx: %w", err)
	}

	if err := client.SendTransaction(ctx, signed); err != nil {
		return AnchorReceipt{}, fmt.Errorf("send tx: %w", err)
	}

	return AnchorReceipt{
		TxHash:      signed.Hash().Hex(),
		Timestamp:   time.Unix(int64(header.Time), 0).UTC(),
		BlockNumber: header.Number.Uint64(),
	}, nil
}
