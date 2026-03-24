package verify

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/user/aletheia-api/internal/blockchain"
	"github.com/user/aletheia-api/internal/hash"
	"github.com/user/aletheia-api/internal/merkle"
	"github.com/user/aletheia-api/internal/models"
	"github.com/user/aletheia-api/internal/storage"
	"github.com/user/aletheia-api/internal/trust"
)

type Service struct {
	storage    storage.ImageStorage
	bigchain   storage.MetadataStorage
	blockchain blockchain.Blockchain
	trust      *trust.ScoringEngine
}

func NewService(s storage.ImageStorage, b storage.MetadataStorage, bc blockchain.Blockchain, t *trust.ScoringEngine) *Service {
	return &Service{
		storage:    s,
		bigchain:   b,
		blockchain: bc,
		trust:      t,
	}
}

func (s *Service) ProcessUpload(ctx context.Context, imgReader io.Reader) (*models.ImageProof, error) {
	// 1. Read image data
	data, err := io.ReadAll(imgReader)
	if err != nil {
		return nil, err
	}

	// 2. Generate SHA-256
	sha256Hash, err := hash.GenerateSHA256(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// 3. Decode image for pHash and Merkle Tree
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 4. Generate pHash
	phash, err := hash.GeneratePHash(img)
	if err != nil {
		return nil, err
	}

	// 5. Generate Merkle Tree & Root
	mt, err := merkle.NewMerkleTree(img)
	if err != nil {
		return nil, err
	}
	merkleRoot := mt.GetRoot()

	// 6. Upload to IPFS
	cid, err := s.storage.UploadImage(ctx, data)
	if err != nil {
		return nil, err
	}

	// 7. Store proof on Blockchain
	timestamp := time.Now().Unix()
	txHash, err := s.blockchain.StoreProof(ctx, sha256Hash, cid, merkleRoot, timestamp)
	if err != nil {
		// Log error but proceed? For MVP we might want to know.
		fmt.Printf("Blockchain storage failed: %v\n", err)
	}

	proof := &models.ImageProof{
		ProofID:    uuid.New().String(),
		SHA256:     sha256Hash,
		PHash:      phash,
		MerkleRoot: merkleRoot,
		IPFSCID:    cid,
		TxHash:     txHash,
		Timestamp:  time.Unix(timestamp, 0),
		Status:     "COMPLETED",
	}

	// 8. Store metadata in BigchainDB
	err = s.bigchain.StoreMetadata(ctx, proof)
	if err != nil {
		fmt.Printf("BigchainDB storage failed: %v\n", err)
	}

	return proof, nil
}

func (s *Service) VerifyImage(ctx context.Context, imgReader io.Reader) (*models.VerificationResult, error) {
	// In a real scenario, we would fetch the stored proof from BigchainDB or Blockchain
	// For MVP, we simulate the verification logic using the Trust Score Engine.

	// 1. Process the incoming image
	data, err := io.ReadAll(imgReader)
	if err != nil {
		return nil, err
	}

	sha256Hash, _ := hash.GenerateSHA256(bytes.NewReader(data))
	img, _, _ := image.Decode(bytes.NewReader(data))
	phash, _ := hash.GeneratePHash(img)

	// In a real scenario, we use these to query the database
	_ = sha256Hash
	_ = phash

	// Placeholder: In a real app, we'd query BigchainDB for a matching SHA-256 or similar pHash
	// matchedProof := queryDatabase(sha256Hash)

	// Simulating a match for demonstration
	sha256Match := true
	phashSimilarity := 100 // 0-100 scale
	metadataIntegrity := true
	witnessNodes := false

	score := s.trust.CalculateScore(sha256Match, phashSimilarity, metadataIntegrity, witnessNodes)
	verdict := s.trust.GetVerdict(score, sha256Match)

	return &models.VerificationResult{
		TrustScore:      score,
		Verdict:         verdict,
		MatchedProofID:  "dummy-id",
		SimilarityScore: phashSimilarity,
	}, nil
}
