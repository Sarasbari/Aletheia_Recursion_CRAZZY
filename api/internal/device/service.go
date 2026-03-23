package device

import (
	"errors"
	"strings"
	"time"
)

type Proof struct {
	PublicKey        string
	Signature        string
	CaptureTimestamp time.Time
	ParentHash       string
}

func ParseCaptureTimestamp(v string) (time.Time, error) {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return time.Time{}, errors.New("capture timestamp is required")
	}
	t, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func ValidateUploadProof(publicKey, signature string, captureTs time.Time) error {
	if strings.TrimSpace(publicKey) == "" {
		return errors.New("device public key is required")
	}
	if strings.TrimSpace(signature) == "" {
		return errors.New("device signature is required")
	}
	if captureTs.IsZero() {
		return errors.New("capture timestamp is required")
	}
	return nil
}
