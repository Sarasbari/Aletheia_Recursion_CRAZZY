package verify

import (
	"context"
	"errors"
	"strings"
	"time"

	contextcheck "aletheia-api/internal/context"
	"aletheia-api/internal/report"
	"aletheia-api/internal/security"
	"aletheia-api/internal/signature"
	"aletheia-api/internal/tamper"
	"aletheia-api/internal/timestamp"
	"aletheia-api/internal/hash"
	"aletheia-api/internal/merkle"
	"aletheia-api/internal/models"
	"aletheia-api/internal/storage"
	"aletheia-api/internal/trust"
)

type Service struct {
	repo   storage.MetadataRepository
	signer signature.Service
}

func NewService(repo storage.MetadataRepository, signer signature.Service) *Service {
	return &Service{repo: repo, signer: signer}
}

func (s *Service) VerifyImage(ctx context.Context, imageBytes []byte, opts models.VerifyOptions) (models.VerifyResponse, error) {
	if len(imageBytes) == 0 {
		return models.VerifyResponse{}, errors.New("empty image")
	}
	sourceType := strings.ToUpper(strings.TrimSpace(opts.SourceType))
	if sourceType == "" {
		sourceType = "UPLOAD"
	}

	sha := hash.SHA256Hex(imageBytes)
	ph, err := hash.ComputePHash(imageBytes)
	if err != nil {
		return models.VerifyResponse{}, err
	}
	merkleRoot, _, err := merkle.BuildMerkleRoot(imageBytes, 16)
	if err != nil {
		return models.VerifyResponse{}, err
	}

	exact, err := s.repo.FindBySHA256(ctx, sha)
	if err != nil {
		return models.VerifyResponse{}, err
	}

	if exact != nil {
		similarity, simErr := hash.SimilarityScore(ph, exact.PHash)
		if simErr != nil {
			similarity = 0
		}
		metaIntegrity := exact.MerkleRoot == merkleRoot

		verifyKey := exact.PublicKey
		verifySig := exact.Signature
		if opts.DevicePublicKey != "" {
			verifyKey = opts.DevicePublicKey
		}
		if opts.DeviceSignature != "" {
			verifySig = opts.DeviceSignature
		}
		signatureValid, _ := s.signer.VerifyHashHex(sha, verifySig, verifyKey)

		replaySuspicious := security.IsReplaySuspicious(exact, opts.CaptureTimestamp)
		contextSuspicious, contextWarnings := contextcheck.EvaluateMismatch(*exact, contextcheck.VerifyInput{
			CaptureTimestamp: opts.CaptureTimestamp,
			Location:         opts.Location,
			DeviceInfo:       opts.DeviceInfo,
		})
		timestampValid, tsWarnings := timestamp.Validate(*exact, time.Now().UTC())
		fraudFlags := models.FraudFlags{
			Duplicate:  false,
			Spam:       false,
			Replay:     replaySuspicious || contextSuspicious,
			Suspicious: replaySuspicious || contextSuspicious,
		}
		warnings := make([]string, 0, len(tsWarnings)+len(contextWarnings)+1)
		warnings = append(warnings, tsWarnings...)
		warnings = append(warnings, contextWarnings...)
		if replaySuspicious {
			warnings = append(warnings, "possible replay detected from capture timestamp mismatch")
		}

		replaySafe := !replaySuspicious
		similarityFloat := float64(similarity) / 100.0

		score := trust.CalculateExplainableScore(true, signatureValid, similarityFloat, metaIntegrity, replaySafe)
		b := trust.BuildBreakdown(true, signatureValid, similarityFloat, metaIntegrity, replaySafe)

		regions := tamper.PlaceholderRegions(true, metaIntegrity)
		mapped := make([]models.TamperedRegion, 0, len(regions))
		for _, r := range regions {
			mapped = append(mapped, models.TamperedRegion{X: r.X, Y: r.Y})
		}

		verdict := "AUTHENTIC"
		if replaySuspicious || contextSuspicious {
			verdict = "UNKNOWN"
		} else if !metaIntegrity {
			verdict = "TAMPERED"
		}
		proofReport := report.BuildProofReport(exact, verdict, score, similarityFloat, b.SignatureStatus, fraudFlags)
		proofReport.Breakdown.Hash = b.HashStatus

		return models.VerifyResponse{
			Verdict:         verdict,
			TrustScore:      score,
			TimestampValid:  timestampValid,
			SourceType:      sourceType,
			ProofReport:     proofReport,
			Warnings:        warnings,
			Suspicious:      fraudFlags.Suspicious,
			MatchedProofID:  exact.ProofID,
			SimilarityScore: similarity,
			Breakdown: models.VerifyBreakdown{
				HashMatch:      b.HashMatch,
				SignatureValid: b.SignatureValid,
				Similarity:     b.Similarity,
				MetadataValid:  b.MetadataValid,
				ReplaySafe:     b.ReplaySafe,
				Hash:           b.HashStatus,
				Signature:      b.SignatureStatus,
			},
			Device: models.DeviceVerification{
				Verified:  signatureValid,
				PublicKey: verifyKey,
			},
			TamperedRegions:  mapped,
			ReplaySuspicious: replaySuspicious,
		}, nil
	}

	all, err := s.repo.List(ctx)
	if err != nil {
		return models.VerifyResponse{}, err
	}
	if len(all) == 0 {
		flags := models.FraudFlags{Duplicate: false, Spam: false, Replay: false, Suspicious: false}
		unknownReport := report.BuildProofReport(nil, "UNKNOWN", 0, 0, "INVALID", flags)
		return models.VerifyResponse{
			Verdict:         "UNKNOWN",
			TrustScore:      0,
			TimestampValid:  false,
			SourceType:      sourceType,
			ProofReport:     unknownReport,
			Warnings:        []string{"no anchored proof found for this media"},
			Suspicious:      false,
			MatchedProofID:  "",
			SimilarityScore: 0,
			Breakdown: models.VerifyBreakdown{
				HashMatch:      false,
				SignatureValid: false,
				Similarity:     0,
				MetadataValid:  false,
				ReplaySafe:     true,
				Hash:           "MISMATCH",
				Signature:      "INVALID",
			},
			Device: models.DeviceVerification{Verified: false, PublicKey: ""},
			TamperedRegions:  []models.TamperedRegion{{X: 10, Y: 20}, {X: 11, Y: 20}},
			ReplaySuspicious: false,
		}, nil
	}

	best := all[0]
	bestSimilarity := 0
	for _, rec := range all {
		score, simErr := hash.SimilarityScore(ph, rec.PHash)
		if simErr != nil {
			continue
		}
		if score > bestSimilarity {
			bestSimilarity = score
			best = rec
		}
	}

	verdict := "TAMPERED"
	if bestSimilarity >= 80 {
		verdict = "SIMILAR"
	}

	similarityFloat := float64(bestSimilarity) / 100.0
	score := trust.CalculateExplainableScore(false, false, similarityFloat, false, true)
	b := trust.BuildBreakdown(false, false, similarityFloat, false, true)
	flags := models.FraudFlags{Duplicate: false, Spam: false, Replay: false, Suspicious: false}
	proofReport := report.BuildProofReport(&best, verdict, score, similarityFloat, b.SignatureStatus, flags)
	proofReport.Breakdown.Hash = b.HashStatus
	regions := tamper.PlaceholderRegions(false, false)
	mapped := make([]models.TamperedRegion, 0, len(regions))
	for _, r := range regions {
		mapped = append(mapped, models.TamperedRegion{X: r.X, Y: r.Y})
	}

	return models.VerifyResponse{
		Verdict:         verdict,
		TrustScore:      score,
		TimestampValid:  false,
		SourceType:      sourceType,
		ProofReport:     proofReport,
		Warnings:        []string{"exact hash not found; best similarity match returned"},
		Suspicious:      false,
		MatchedProofID:  best.ProofID,
		SimilarityScore: bestSimilarity,
		Breakdown: models.VerifyBreakdown{
			HashMatch:      b.HashMatch,
			SignatureValid: b.SignatureValid,
			Similarity:     b.Similarity,
			MetadataValid:  b.MetadataValid,
			ReplaySafe:     b.ReplaySafe,
			Hash:           b.HashStatus,
			Signature:      b.SignatureStatus,
		},
		Device: models.DeviceVerification{Verified: false, PublicKey: best.PublicKey},
		TamperedRegions:  mapped,
		ReplaySuspicious: false,
	}, nil
}
