package fraud

import (
	"strings"
	"sync"
	"time"

	"aletheia-api/internal/models"
	"aletheia-api/internal/security"
)

type Detector struct {
	mu                sync.Mutex
	uploadWindow      time.Duration
	maxUploadsPerSpan int
	uploadsByActor    map[string][]time.Time
}

func NewDetector(uploadWindow time.Duration, maxUploadsPerSpan int) *Detector {
	if uploadWindow <= 0 {
		uploadWindow = 10 * time.Minute
	}
	if maxUploadsPerSpan <= 0 {
		maxUploadsPerSpan = 20
	}
	return &Detector{
		uploadWindow:      uploadWindow,
		maxUploadsPerSpan: maxUploadsPerSpan,
		uploadsByActor:    make(map[string][]time.Time),
	}
}

func (d *Detector) RegisterUploadAndDetectSpam(actor string, now time.Time) bool {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		return false
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := now.Add(-d.uploadWindow)
	entries := d.uploadsByActor[actor]
	filtered := make([]time.Time, 0, len(entries)+1)
	for _, ts := range entries {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	filtered = append(filtered, now)
	d.uploadsByActor[actor] = filtered

	return len(filtered) > d.maxUploadsPerSpan
}

func EvaluateUpload(existing *models.ProofRecord, captureTimestamp time.Time, location, deviceInfo string, spam bool) (models.FraudFlags, []string) {
	flags := models.FraudFlags{Duplicate: existing != nil, Spam: spam}
	warnings := make([]string, 0, 3)

	if flags.Duplicate {
		warnings = append(warnings, "duplicate upload detected for existing hash")
	}
	if flags.Spam {
		warnings = append(warnings, "high upload rate detected for actor")
	}

	if existing != nil {
		replayMismatch := security.IsReplaySuspicious(existing, captureTimestamp)
		contextMismatch := false
		if strings.TrimSpace(location) != "" && strings.TrimSpace(existing.Location) != "" && !strings.EqualFold(strings.TrimSpace(location), strings.TrimSpace(existing.Location)) {
			contextMismatch = true
		}
		if strings.TrimSpace(deviceInfo) != "" && strings.TrimSpace(existing.DeviceInfo) != "" && !strings.EqualFold(strings.TrimSpace(deviceInfo), strings.TrimSpace(existing.DeviceInfo)) {
			contextMismatch = true
		}
		flags.Replay = replayMismatch || contextMismatch
		if flags.Replay {
			warnings = append(warnings, "possible replay attack with altered metadata")
		}
	}

	flags.Suspicious = flags.Spam || flags.Replay
	return flags, warnings
}
