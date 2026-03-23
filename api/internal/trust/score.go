package trust

import "aletheia-api/internal/models"

func CalculateScore(input models.TrustInput) int {
	score := 0
	if input.SHA256Match {
		score += 35
	}
	if input.PHashSimilarity {
		score += 20
	}
	if input.MetadataIntegrity {
		score += 10
	}
	if input.WitnessConsensus {
		score += 15
	}
	if score > 100 {
		return 100
	}
	if score < 0 {
		return 0
	}
	return score
}
