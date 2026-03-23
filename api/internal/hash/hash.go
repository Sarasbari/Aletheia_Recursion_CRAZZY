package hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strconv"
	"strings"
)

func SHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func SHA256File(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func ComputePHash(imageBytes []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", err
	}

	gray := resizeToGray32(img)
	dct := dct2D(gray)

	coeffs := make([]float64, 0, 64)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if x == 0 && y == 0 {
				continue
			}
			coeffs = append(coeffs, dct[y][x])
		}
	}

	if len(coeffs) == 0 {
		return "", errors.New("unable to compute pHash coefficients")
	}

	median := median(coeffs)
	var hash uint64
	for i, c := range coeffs {
		if c > median {
			hash |= 1 << uint(i)
		}
	}

	return fmt.Sprintf("%016x", hash), nil
}

func HammingDistanceHex(aHex, bHex string) (int, error) {
	a, err := strconv.ParseUint(strings.TrimSpace(aHex), 16, 64)
	if err != nil {
		return 0, err
	}
	b, err := strconv.ParseUint(strings.TrimSpace(bHex), 16, 64)
	if err != nil {
		return 0, err
	}
	x := a ^ b
	count := 0
	for x > 0 {
		x &= x - 1
		count++
	}
	return count, nil
}

func SimilarityScore(aHex, bHex string) (int, error) {
	d, err := HammingDistanceHex(aHex, bHex)
	if err != nil {
		return 0, err
	}
	if d < 0 {
		d = 0
	}
	if d > 64 {
		d = 64
	}
	return int(math.Round((1 - float64(d)/64.0) * 100)), nil
}

func resizeToGray32(img image.Image) [32][32]float64 {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	var out [32][32]float64

	for y := 0; y < 32; y++ {
		sy := bounds.Min.Y + int(float64(y)*float64(h)/32.0)
		if sy >= bounds.Max.Y {
			sy = bounds.Max.Y - 1
		}
		for x := 0; x < 32; x++ {
			sx := bounds.Min.X + int(float64(x)*float64(w)/32.0)
			if sx >= bounds.Max.X {
				sx = bounds.Max.X - 1
			}
			r, g, b, _ := img.At(sx, sy).RGBA()
			gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
			out[y][x] = gray
		}
	}
	return out
}

func dct2D(input [32][32]float64) [32][32]float64 {
	var out [32][32]float64
	for v := 0; v < 32; v++ {
		for u := 0; u < 32; u++ {
			sum := 0.0
			for y := 0; y < 32; y++ {
				for x := 0; x < 32; x++ {
					sum += input[y][x] *
						math.Cos((float64(2*x+1)*float64(u)*math.Pi)/64.0) *
						math.Cos((float64(2*y+1)*float64(v)*math.Pi)/64.0)
				}
			}
			tauU := 1.0
			if u == 0 {
				tauU = 1 / math.Sqrt2
			}
			tauV := 1.0
			if v == 0 {
				tauV = 1 / math.Sqrt2
			}
			out[v][u] = 0.25 * tauU * tauV * sum
		}
	}
	return out
}

func median(vals []float64) float64 {
	sorted := append([]float64(nil), vals...)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] < sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}
