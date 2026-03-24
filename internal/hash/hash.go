package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/corona10/goimagehash"
)

// GenerateSHA256 returns the SHA-256 hash of the given reader.
func GenerateSHA256(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// GeneratePHash returns the perceptual hash of the given image.
func GeneratePHash(img image.Image) (string, error) {
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", err
	}
	return hash.ToString(), nil
}

// ComparePHash returns the distance between two perceptual hashes.
func ComparePHash(hash1Str, hash2Str string) (int, error) {
	h1, err := goimagehash.ImageHashFromString(hash1Str)
	if err != nil {
		return 0, err
	}
	h2, err := goimagehash.ImageHashFromString(hash2Str)
	if err != nil {
		return 0, err
	}
	distance, err := h1.Distance(h2)
	if err != nil {
		return 0, err
	}
	return distance, nil
}
