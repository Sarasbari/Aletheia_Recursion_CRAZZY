package clients

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"aletheia-api/shared/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PolygonClient struct {
	rpcURL          string
	privateKeyHex   string
	contractAddress string
	chainID         *big.Int
}

func NewPolygonClient(rpcURL, privateKeyHex, contractAddress, chainID string) *PolygonClient {
	cid := big.NewInt(137)
	if parsed, ok := new(big.Int).SetString(chainID, 10); ok {
		cid = parsed
	}
	return &PolygonClient{
		rpcURL:          rpcURL,
		privateKeyHex:   privateKeyHex,
		contractAddress: contractAddress,
		chainID:         cid,
	}
}

func (p *PolygonClient) Anchor(ctx context.Context, record models.ProvenRecord) (string, error) {
	if p.rpcURL == "" || p.privateKeyHex == "" || p.contractAddress == "" {
		body, _ := json.Marshal(record)
		sum := sha256.Sum256(body)
		return "mock-" + hex.EncodeToString(sum[:]), nil
	}

	client, err := ethclient.DialContext(ctx, p.rpcURL)
	if err != nil {
		return "", err
	}
	defer client.Close()

	pk, err := crypto.HexToECDSA(p.privateKeyHex)
	if err != nil {
		return "", err
	}
	from := crypto.PubkeyToAddress(pk.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}
	payload, _ := json.Marshal(record)
	tx := types.NewTransaction(nonce, common.HexToAddress(p.contractAddress), big.NewInt(0), 300000, gasPrice, payload)
	signed, err := types.SignTx(tx, types.NewEIP155Signer(p.chainID), pk)
	if err != nil {
		return "", err
	}
	if err := client.SendTransaction(ctx, signed); err != nil {
		return "", fmt.Errorf("send polygon tx: %w", err)
	}
	return signed.Hash().Hex(), nil
}
