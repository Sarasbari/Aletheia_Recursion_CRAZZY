package context

import (
	"strings"
	"time"

	"aletheia-api/internal/models"
)

const allowedCaptureDelta = 5 * time.Second

type VerifyInput struct {
	CaptureTimestamp time.Time
	Location         string
	DeviceInfo       string
}

func EvaluateMismatch(record models.ProofRecord, input VerifyInput) (bool, []string) {
	warnings := make([]string, 0, 3)
	suspicious := false

	if !input.CaptureTimestamp.IsZero() && !record.CaptureTimestamp.IsZero() {
		delta := record.CaptureTimestamp.Sub(input.CaptureTimestamp)
		if delta < 0 {
			delta = -delta
		}
		if delta > allowedCaptureDelta {
			suspicious = true
			warnings = append(warnings, "capture time mismatch with anchored proof")
		}
	}

	if normalize(record.Location) != "" && normalize(input.Location) != "" && normalize(record.Location) != normalize(input.Location) {
		suspicious = true
		warnings = append(warnings, "location mismatch with anchored context")
	}

	if normalize(record.DeviceInfo) != "" && normalize(input.DeviceInfo) != "" && normalize(record.DeviceInfo) != normalize(input.DeviceInfo) {
		suspicious = true
		warnings = append(warnings, "device info mismatch with anchored context")
	}

	return suspicious, warnings
}

func normalize(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}
