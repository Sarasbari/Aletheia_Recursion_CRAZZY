package security

import (
	"time"

	"aletheia-api/internal/models"
)

const replayTolerance = 5 * time.Second

func IsReplaySuspicious(existing *models.ProofRecord, newCaptureTimestamp time.Time) bool {
	if existing == nil || newCaptureTimestamp.IsZero() || existing.CaptureTimestamp.IsZero() {
		return false
	}
	delta := existing.CaptureTimestamp.Sub(newCaptureTimestamp)
	if delta < 0 {
		delta = -delta
	}
	return delta > replayTolerance
}
