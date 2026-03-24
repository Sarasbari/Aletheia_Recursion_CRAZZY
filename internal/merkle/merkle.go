package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	"image/draw"
)

// MerkleTree represents a binary Merkle Tree.
type MerkleTree struct {
	RootNode *Node
}

// Node represents a node in the Merkle Tree.
type Node struct {
	Hash  string
	Left  *Node
	Right *Node
}

// NewMerkleTree builds a binary Merkle Tree from image blocks.
func NewMerkleTree(img image.Image) (*MerkleTree, error) {
	blocks := splitImageIntoBlocks(img, 16, 16)
	if len(blocks) == 0 {
		return nil, errors.New("no blocks generated from image")
	}

	var nodes []Node
	for _, block := range blocks {
		nodes = append(nodes, Node{Hash: hashBlock(block)})
	}

	return &MerkleTree{RootNode: buildTree(nodes)}, nil
}

func splitImageIntoBlocks(img image.Image, blockWidth, blockHeight int) []image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	var blocks []image.Image

	for y := 0; y < height; y += blockHeight {
		for x := 0; x < width; x += blockWidth {
			rect := image.Rect(x, y, x+blockWidth, y+blockHeight)
			// Ensure we don't go out of bounds
			if rect.Max.X > width {
				rect.Max.X = width
			}
			if rect.Max.Y > height {
				rect.Max.Y = height
			}

			// Create a new RGBA image for the block
			block := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
			draw.Draw(block, block.Bounds(), img, rect.Min, draw.Src)
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func hashBlock(img image.Image) string {
	h := sha256.New()
	// Use image pixel data for hashing
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			h.Write([]byte{byte(r), byte(g), byte(b), byte(a)})
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

func buildTree(nodes []Node) *Node {
	if len(nodes) == 1 {
		return &nodes[0]
	}

	var newLevel []Node
	for i := 0; i < len(nodes); i += 2 {
		var left, right *Node
		left = &nodes[i]
		if i+1 < len(nodes) {
			right = &nodes[i+1]
		} else {
			// Duplicate the last node if the number of nodes is odd
			right = &nodes[i]
		}

		combinedHash := sha256.Sum256([]byte(left.Hash + right.Hash))
		newLevel = append(newLevel, Node{
			Hash:  hex.EncodeToString(combinedHash[:]),
			Left:  left,
			Right: right,
		})
	}

	return buildTree(newLevel)
}

// GetRoot returns the hex string of the Merkle Root.
func (mt *MerkleTree) GetRoot() string {
	if mt.RootNode == nil {
		return ""
	}
	return mt.RootNode.Hash
}
