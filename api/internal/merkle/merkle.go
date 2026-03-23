package merkle

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type ProofNode struct {
	Hash string `json:"hash"`
	Left bool   `json:"left"`
}

func BuildMerkleRoot(imageBytes []byte, blockSize int) (string, []string, error) {
	if blockSize <= 0 {
		return "", nil, errors.New("block size must be positive")
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", nil, err
	}

	leaves := buildLeafHashes(img, blockSize)
	if len(leaves) == 0 {
		return "", nil, errors.New("no leaves generated")
	}

	root := buildRoot(leaves)
	return root, leaves, nil
}

func VerifyMerkleProof(leafHash string, proof []ProofNode, rootHash string) bool {
	current, err := hex.DecodeString(leafHash)
	if err != nil {
		return false
	}

	for _, node := range proof {
		sibling, err := hex.DecodeString(node.Hash)
		if err != nil {
			return false
		}
		if node.Left {
			current = hashPair(sibling, current)
		} else {
			current = hashPair(current, sibling)
		}
	}

	return hex.EncodeToString(current) == rootHash
}

func verifyMerkleProof(leafHash string, proof []ProofNode, rootHash string) bool {
	return VerifyMerkleProof(leafHash, proof, rootHash)
}

func buildLeafHashes(img image.Image, blockSize int) []string {
	b := img.Bounds()
	leaves := make([]string, 0)

	for y := b.Min.Y; y < b.Max.Y; y += blockSize {
		for x := b.Min.X; x < b.Max.X; x += blockSize {
			var block []byte
			block = append(block, byte(x&0xff), byte(y&0xff))

			yMax := min(y+blockSize, b.Max.Y)
			xMax := min(x+blockSize, b.Max.X)
			for by := y; by < yMax; by++ {
				for bx := x; bx < xMax; bx++ {
					r, g, bl, a := img.At(bx, by).RGBA()
					block = append(block,
						byte(r>>8), byte(g>>8), byte(bl>>8), byte(a>>8),
					)
				}
			}

			sum := sha256.Sum256(block)
			leaves = append(leaves, hex.EncodeToString(sum[:]))
		}
	}
	return leaves
}

func buildRoot(leaves []string) string {
	layer := make([][]byte, 0, len(leaves))
	for _, leaf := range leaves {
		decoded, _ := hex.DecodeString(leaf)
		layer = append(layer, decoded)
	}

	for len(layer) > 1 {
		next := make([][]byte, 0, (len(layer)+1)/2)
		for i := 0; i < len(layer); i += 2 {
			left := layer[i]
			right := left
			if i+1 < len(layer) {
				right = layer[i+1]
			}
			next = append(next, hashPair(left, right))
		}
		layer = next
	}
	return hex.EncodeToString(layer[0])
}

func hashPair(a, b []byte) []byte {
	payload := append(append([]byte{}, a...), b...)
	sum := sha256.Sum256(payload)
	return sum[:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

