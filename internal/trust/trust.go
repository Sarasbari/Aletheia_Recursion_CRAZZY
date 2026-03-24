package trust

import (
	"github.com/user/aletheia-api/internal/models"
)

// ScoringEngine handles the calculation of trust scores.
type ScoringEngine struct{}

func NewScoringEngine() *ScoringEngine {
	return &ScoringEngine{}
}

// CalculateScore computes the trust score based on specified criteria.
func (e *ScoringEngine) CalculateScore(sha256Match bool, phashSimilarity int, metadataIntegrity bool, witnessNodes bool) int {
	score := 0

	// SHA256 match → +35
	if sha256Match {
		score += 35
	}

	// pHash similarity → +20 (Assume distance < 10 is similar)
	if phashSimilarity > 90 { // Similarity percentage or inverse distance
		score += 20
	} else if phashSimilarity > 70 {
		score += 10
	}

	// metadata integrity → +10
	if metadataIntegrity {
		score += 10
	}

	// (optional placeholder for witness nodes) → +15
	if witnessNodes {
		score += 15
	}

	// Add more criteria to reach 100 if necessary
	// Current max: 35 + 20 + 10 + 15 = 80.
	// We can scale it or add 20 for "on-chain verification"
	if sha256Match && metadataIntegrity {
		score += 20 // Bonus for high integrity
	}

	if score > 100 {
		score = 100
	}

	return score
}

// GetVerdict returns the verdict based on the trust score.
func (e *ScoringEngine) GetVerdict(score int, sha256Match bool) string {
	if sha256Match && score >= 90 {
		return models.VerdictAuthentic
	}
	if score >= 50 {
		return models.VerdictSimilar
	}
	if score > 0 {
		return models.VerdictTampered
	}
	return models.VerdictUnknown
}
