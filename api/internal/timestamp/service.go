package timestamp

import (
	"fmt"
	"time"

	"aletheia-api/internal/models"
)

const maxClockSkew = 2 * time.Minute

type Integrity struct {
	Hash        string    `json:"hash"`
	Timestamp   time.Time `json:"timestamp"`
	BlockNumber uint64    `json:"blockNumber"`
}

func FromRecord(record models.ProofRecord) Integrity {
	return Integrity{
		Hash:        record.SHA256,
		Timestamp:   record.Timestamp,
		BlockNumber: record.BlockNumber,
	}
}

func Validate(record models.ProofRecord, now time.Time) (bool, []string) {
	warnings := make([]string, 0, 3)

	if record.Timestamp.IsZero() {
		warnings = append(warnings, "missing anchored timestamp")
		return false, warnings
	}
	if record.BlockNumber == 0 {
		warnings = append(warnings, "missing block number for anchored proof")
		return false, warnings
	}
	if record.Timestamp.After(now.Add(maxClockSkew)) {
		warnings = append(warnings, "anchored timestamp is in the future")
		return false, warnings
	}
	if !record.CaptureTimestamp.IsZero() && record.Timestamp.Before(record.CaptureTimestamp.Add(-maxClockSkew)) {
		warnings = append(warnings, fmt.Sprintf("anchored timestamp (%s) predates capture time (%s)", record.Timestamp.Format(time.RFC3339), record.CaptureTimestamp.Format(time.RFC3339)))
		return false, warnings
	}

	return true, warnings
}
