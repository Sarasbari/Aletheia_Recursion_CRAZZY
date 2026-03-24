package merkle

import (
	"image"
	"image/color"
	"testing"
)

func TestMerkleTree(t *testing.T) {
	// Create a dummy 32x32 image
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 100, 255})
		}
	}

	mt, err := NewMerkleTree(img)
	if err != nil {
		t.Fatalf("Failed to create Merkle Tree: %v", err)
	}

	root := mt.GetRoot()
	if root == "" {
		t.Fatal("Merkle Root should not be empty")
	}

	t.Logf("Generated Merkle Root: %s", root)

	// Verify consistency
	mt2, _ := NewMerkleTree(img)
	if mt2.GetRoot() != root {
		t.Error("Merkle Root should be consistent for the same image")
	}

	// Change one pixel
	img.Set(0, 0, color.RGBA{255, 255, 255, 255})
	mt3, _ := NewMerkleTree(img)
	if mt3.GetRoot() == root {
		t.Error("Merkle Root should change when the image is modified")
	}
}
