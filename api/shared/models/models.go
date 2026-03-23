package models

import "time"

type UploadJob struct {
	JobID       string    `json:"jobId"`
	FileName    string    `json:"fileName"`
	ImageBase64 string    `json:"imageBase64"`
	UploadedAt  time.Time `json:"uploadedAt"`
}

type HashResult struct {
	JobID       string    `json:"jobId"`
	FileName    string    `json:"fileName"`
	ImageBase64 string    `json:"imageBase64"`
	SHA256      string    `json:"sha256"`
	PHash       string    `json:"phash"`
	ProcessedAt time.Time `json:"processedAt"`
}

type MerkleResult struct {
	JobID       string    `json:"jobId"`
	FileName    string    `json:"fileName"`
	ImageBase64 string    `json:"imageBase64"`
	SHA256      string    `json:"sha256"`
	PHash       string    `json:"phash"`
	MerkleRoot  string    `json:"merkleRoot"`
	ProcessedAt time.Time `json:"processedAt"`
}

type AIResult struct {
	JobID         string    `json:"jobId"`
	AIProbability float64   `json:"aiProbability"`
	ProcessedAt   time.Time `json:"processedAt"`
}

type ProvenRecord struct {
	JobID       string    `json:"jobId"`
	SHA256      string    `json:"sha256"`
	PHash       string    `json:"phash"`
	MerkleRoot  string    `json:"merkleRoot"`
	IPFSCID     string    `json:"ipfsCID"`
	TxHash      string    `json:"txHash"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
	FileName    string    `json:"fileName"`
	ImageBase64 string    `json:"imageBase64,omitempty"`
}

type VerifyResponse struct {
	TrustScore      int    `json:"trustScore"`
	Verdict         string `json:"verdict"`
	MatchedProofID  string `json:"matchedProofId"`
	SimilarityScore int    `json:"similarityScore"`
}
