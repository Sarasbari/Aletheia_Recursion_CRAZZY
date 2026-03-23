package trust

type Breakdown struct {
	HashMatch      bool    `json:"hashMatch"`
	SignatureValid bool    `json:"signatureValid"`
	Similarity     float64 `json:"similarity"`
	MetadataValid  bool    `json:"metadataValid"`
	ReplaySafe     bool    `json:"replaySafe"`
	HashStatus     string  `json:"hash"`
	SignatureStatus string `json:"signature"`
}

func CalculateExplainableScore(hashMatch, signatureValid bool, similarity float64, metadataValid, replaySafe bool) int {
	score := 0
	if hashMatch {
		score += 35
	}
	if similarity >= 0.80 {
		score += 20
	}
	if metadataValid {
		score += 10
	}
	if signatureValid {
		score += 15
	}
	if !replaySafe {
		score -= 20
	}
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

func BuildBreakdown(hashMatch, signatureValid bool, similarity float64, metadataValid, replaySafe bool) Breakdown {
	hashStatus := "MISMATCH"
	if hashMatch {
		hashStatus = "MATCH"
	}
	sigStatus := "INVALID"
	if signatureValid {
		sigStatus = "VALID"
	}
	return Breakdown{
		HashMatch:       hashMatch,
		SignatureValid:  signatureValid,
		Similarity:      similarity,
		MetadataValid:   metadataValid,
		ReplaySafe:      replaySafe,
		HashStatus:      hashStatus,
		SignatureStatus: sigStatus,
	}
}
