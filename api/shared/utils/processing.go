package utils

import (
	"encoding/base64"

	"aletheia-api/internal/hash"
	"aletheia-api/internal/merkle"
	"aletheia-api/internal/models"
	"aletheia-api/internal/trust"
)

func DecodeBase64Image(v string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(v)
}

func EncodeBase64Image(v []byte) string {
	return base64.StdEncoding.EncodeToString(v)
}

func ComputeSHA256AndPHash(image []byte) (string, string, error) {
	sha := hash.SHA256Hex(image)
	ph, err := hash.ComputePHash(image)
	if err != nil {
		return "", "", err
	}
	return sha, ph, nil
}

func ComputeMerkleRoot(image []byte) (string, error) {
	root, _, err := merkle.BuildMerkleRoot(image, 16)
	if err != nil {
		return "", err
	}
	return root, nil
}

func SimilarityScore(a, b string) (int, error) {
	return hash.SimilarityScore(a, b)
}

func ComputeTrustScore(shaMatch, pHashSimilarity, metadataIntegrity, witness bool) int {
	return trust.CalculateScore(models.TrustInput{
		SHA256Match:       shaMatch,
		PHashSimilarity:   pHashSimilarity,
		MetadataIntegrity: metadataIntegrity,
		WitnessConsensus:  witness,
	})
}
