package signature

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type Service interface {
	SignHashHex(hashHex, privateKeyHex string) (signatureHex string, publicKeyHex string, err error)
	VerifyHashHex(hashHex, signatureHex, publicKeyHex string) (bool, error)
}

type Secp256k1Service struct{}

func NewSecp256k1Service() *Secp256k1Service {
	return &Secp256k1Service{}
}

func (s *Secp256k1Service) SignHashHex(hashHex, privateKeyHex string) (string, string, error) {
	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		return "", "", fmt.Errorf("decode hash hex: %w", err)
	}
	pk, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", "", fmt.Errorf("decode private key: %w", err)
	}
	sig, err := crypto.Sign(hashBytes, pk)
	if err != nil {
		return "", "", fmt.Errorf("sign hash: %w", err)
	}
	pubBytes := crypto.FromECDSAPub(&pk.PublicKey)
	return hex.EncodeToString(sig), hex.EncodeToString(pubBytes), nil
}

func (s *Secp256k1Service) VerifyHashHex(hashHex, signatureHex, publicKeyHex string) (bool, error) {
	if hashHex == "" || signatureHex == "" || publicKeyHex == "" {
		return false, errors.New("hash, signature and publicKey are required")
	}

	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		return false, fmt.Errorf("decode hash hex: %w", err)
	}
	sig, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("decode signature hex: %w", err)
	}
	pubBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("decode public key hex: %w", err)
	}

	pubKey, err := crypto.UnmarshalPubkey(pubBytes)
	if err != nil {
		return false, fmt.Errorf("unmarshal public key: %w", err)
	}

	if len(sig) != 65 {
		return false, errors.New("signature must be 65 bytes")
	}

	verifyBytes := pubBytes
	if len(verifyBytes) == 65 && verifyBytes[0] == 0x04 {
		verifyBytes = verifyBytes[1:]
	}
	verified := crypto.VerifySignature(verifyBytes, hashBytes, sig[:64])
	if verified {
		return true, nil
	}

	// Fallback: recover pubkey and compare.
	recovered, recErr := crypto.SigToPub(hashBytes, sig)
	if recErr != nil {
		return false, nil
	}
	return recovered.X.Cmp(pubKey.X) == 0 && recovered.Y.Cmp(pubKey.Y) == 0, nil
}
