package report

import (
	"time"

	"aletheia-api/internal/models"
)

func BuildProofReport(record *models.ProofRecord, verdict string, trustScore int, similarity float64, signatureState string, flags models.FraudFlags) models.ProofReport {
	report := models.ProofReport{
		ProofID:    "",
		SHA256:     "",
		CID:        "",
		Timestamp:  "",
		BlockNumber: 0,
		Verdict:    verdict,
		TrustScore: trustScore,
		Breakdown: models.ProofReportBreakdown{
			Hash:       "MISMATCH",
			Similarity: similarity,
			Signature:  signatureState,
		},
		Flags: flags,
	}

	if record == nil {
		return report
	}

	report.ProofID = record.ProofID
	report.SHA256 = record.SHA256
	report.CID = record.IPFSCID
	report.BlockNumber = record.BlockNumber
	if !record.Timestamp.IsZero() {
		report.Timestamp = record.Timestamp.UTC().Format(time.RFC3339)
	}
	if verdict == "AUTHENTIC" {
		report.Breakdown.Hash = "MATCH"
	}
	return report
}
