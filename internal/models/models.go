package models

import "time"

// ImageProof represents the cryptographic proof of an image.
type ImageProof struct {
	ProofID    string    `json:"proofId"`
	SHA256     string    `json:"sha256"`
	PHash      string    `json:"phash"`
	MerkleRoot string    `json:"merkleRoot"`
	IPFSCID    string    `json:"ipfsCID"`
	TxHash     string    `json:"txHash"`
	Timestamp  time.Time `json:"timestamp"`
	Status     string    `json:"status"`
}

// VerificationResult represents the outcome of an authenticity check.
type VerificationResult struct {
	TrustScore      int    `json:"trustScore"`
	Verdict         string `json:"verdict"` // AUTHENTIC, TAMPERED, SIMILAR, UNKNOWN
	MatchedProofID  string `json:"matchedProofId"`
	SimilarityScore int    `json:"similarityScore"`
}

const (
	VerdictAuthentic = "AUTHENTIC"
	VerdictTampered  = "TAMPERED"
	VerdictSimilar   = "SIMILAR"
	VerdictUnknown   = "UNKNOWN"
)
