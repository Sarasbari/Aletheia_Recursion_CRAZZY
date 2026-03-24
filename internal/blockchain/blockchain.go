package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Blockchain defines the interface for interacting with Polygon.
type Blockchain interface {
	StoreProof(ctx context.Context, sha256, cid, merkleRoot string, timestamp int64) (string, error)
}

// PolygonService implements Blockchain for Polygon.
type PolygonService struct {
	client          *ethclient.Client
	privateKey      *ecdsa.PrivateKey
	address         common.Address
	contractAddress common.Address
	contractABI     abi.ABI
}

const contractABIJSON = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"string","name":"sha256","type":"string"},{"indexed":false,"internalType":"string","name":"ipfsCID","type":"string"},{"indexed":false,"internalType":"string","name":"merkleRoot","type":"string"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"},{"indexed":true,"internalType":"address","name":"owner","type":"address"}],"name":"ProofStored","type":"event"},{"inputs":[{"internalType":"string","name":"_sha256","type":"string"}],"name":"getProof","outputs":[{"internalType":"string","name":"","type":"string"},{"internalType":"string","name":"","type":"string"},{"internalType":"string","name":"","type":"string"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_sha256","type":"string"},{"internalType":"string","name":"_ipfsCID","type":"string"},{"internalType":"string","name":"_merkleRoot","type":"string"}],"name":"storeProof","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

func NewPolygonService(rawURL, hexPrivateKey, contractAddr string) (*PolygonService, error) {
	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)

	parsedABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	return &PolygonService{
		client:          client,
		privateKey:      privateKey,
		address:         address,
		contractAddress: common.HexToAddress(contractAddr),
		contractABI:     parsedABI,
	}, nil
}

func (s *PolygonService) StoreProof(ctx context.Context, sha256, cid, merkleRoot string, timestamp int64) (string, error) {
	// Pack the transaction data for the storeProof method
	data, err := s.contractABI.Pack("storeProof", sha256, cid, merkleRoot)
	if err != nil {
		return "", fmt.Errorf("failed to pack contract call: %w", err)
	}

	nonce, err := s.client.PendingNonceAt(ctx, s.address)
	if err != nil {
		return "", err
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	// Estimate gas for the contract call
	gasLimit, err := s.client.EstimateGas(ctx, ethereum.CallMsg{
		From: s.address,
		To:   &s.contractAddress,
		Data: data,
	})
	if err != nil {
		gasLimit = uint64(300000) // Fallback gas limit
	}

	tx := types.NewTransaction(nonce, s.contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := s.client.ChainID(ctx)
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), s.privateKey)
	if err != nil {
		return "", err
	}

	err = s.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}
