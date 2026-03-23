package models

import "time"

type ProofRecord struct {
	ProofID    string    `json:"proofId"`
	SHA256     string    `json:"sha256"`
	PHash      string    `json:"phash"`
	MerkleRoot string    `json:"merkleRoot"`
	Signature  string    `json:"signature"`
	PublicKey  string    `json:"publicKey"`
	ParentHash string    `json:"parentHash,omitempty"`
	Location   string    `json:"location,omitempty"`
	DeviceInfo string    `json:"deviceInfo,omitempty"`
	UploaderID string    `json:"uploaderId,omitempty"`
	CaptureTimestamp time.Time `json:"captureTimestamp"`
	IPFSCID    string    `json:"ipfsCID"`
	TxHash     string    `json:"txHash"`
	BlockNumber uint64   `json:"blockNumber"`
	Timestamp  time.Time `json:"timestamp"`
	Status     string    `json:"status"`
}

type UploadResponse struct {
	ProofID    string `json:"proofId"`
	SHA256     string `json:"sha256"`
	PHash      string `json:"phash"`
	MerkleRoot string `json:"merkleRoot"`
	Signature  string `json:"signature"`
	PublicKey  string `json:"publicKey"`
	CaptureTimestamp string `json:"captureTimestamp"`
	IPFSCID    string `json:"ipfsCID"`
	TxHash     string `json:"txHash"`
	BlockNumber uint64 `json:"blockNumber"`
	Status     string `json:"status"`
	Flags      FraudFlags `json:"flags"`
	Warnings   []string `json:"warnings,omitempty"`
}

type VerifyResponse struct {
	Verdict         string `json:"verdict"`
	TrustScore      int    `json:"trustScore"`
	TimestampValid  bool   `json:"timestampValid"`
	SourceType      string `json:"sourceType"`
	ProofReport     ProofReport `json:"proofReport"`
	Warnings        []string `json:"warnings"`
	Suspicious      bool   `json:"suspicious"`
	MatchedProofID  string `json:"matchedProofId"`
	SimilarityScore int    `json:"similarityScore"`
	Breakdown       VerifyBreakdown `json:"breakdown"`
	Device          DeviceVerification `json:"device"`
	TamperedRegions []TamperedRegion `json:"tamperedRegions"`
	ReplaySuspicious bool `json:"replaySuspicious"`
}

type FraudFlags struct {
	Duplicate  bool `json:"duplicate"`
	Spam       bool `json:"spam"`
	Replay     bool `json:"replay"`
	Suspicious bool `json:"suspicious"`
}

type ProofReportBreakdown struct {
	Hash       string  `json:"hash"`
	Similarity float64 `json:"similarity"`
	Signature  string  `json:"signature"`
}

type ProofReport struct {
	ProofID     string               `json:"proofId"`
	SHA256      string               `json:"sha256"`
	CID         string               `json:"CID"`
	Timestamp   string               `json:"timestamp"`
	BlockNumber uint64               `json:"blockNumber"`
	Verdict     string               `json:"verdict"`
	TrustScore  int                  `json:"trustScore"`
	Breakdown   ProofReportBreakdown `json:"breakdown"`
	Flags       FraudFlags           `json:"flags"`
}

type VerifyOptions struct {
	SourceType        string
	CaptureTimestamp  time.Time
	Location          string
	DeviceInfo        string
	DevicePublicKey   string
	DeviceSignature   string
}

type VerifyBreakdown struct {
	HashMatch      bool    `json:"hashMatch"`
	SignatureValid bool    `json:"signatureValid"`
	Similarity     float64 `json:"similarity"`
	MetadataValid  bool    `json:"metadataValid"`
	ReplaySafe     bool    `json:"replaySafe"`
	Hash           string  `json:"hash"`
	Signature      string  `json:"signature"`
}

type DeviceVerification struct {
	Verified  bool   `json:"verified"`
	PublicKey string `json:"publicKey"`
}

type TamperedRegion struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type TrustInput struct {
	SHA256Match       bool
	PHashSimilarity   bool
	MetadataIntegrity bool
	WitnessConsensus  bool
}

type VideoRecord struct {
	ProofID    string    `json:"proofId"`
	VideoHash  string    `json:"videoHash"`
	AudioHash  string    `json:"audioHash"`
	VideoCID   string    `json:"cidVideo"`
	AudioCID   string    `json:"cidAudio"`
	TxHash     string    `json:"txHash"`
	Timestamp  time.Time `json:"timestamp"`
	Status     string    `json:"status"`
	FileName   string    `json:"fileName"`
}

type VideoUploadResponse struct {
	ProofID   string `json:"proofId"`
	VideoHash string `json:"videoHash"`
	AudioHash string `json:"audioHash"`
	CIDVideo  string `json:"cidVideo"`
	CIDAudio  string `json:"cidAudio"`
	TxHash    string `json:"txHash"`
	Status    string `json:"status"`
}

type VideoVerifyResponse struct {
	VideoHash string `json:"videoHash"`
	AudioHash string `json:"audioHash"`
	Verdict   string `json:"verdict"`
	Details   struct {
		Video string `json:"video"`
		Audio string `json:"audio"`
	} `json:"details"`
}
