package video

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aletheia-api/internal/audio"
	"aletheia-api/internal/blockchain"
	"aletheia-api/internal/hash"
	"aletheia-api/internal/models"
	"aletheia-api/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	repo       storage.MetadataRepository
	ipfs       storage.IPFSClient
	bigchain   storage.BigchainClient
	blockchain blockchain.Client
}

func NewService(repo storage.MetadataRepository, ipfs storage.IPFSClient, bigchain storage.BigchainClient, blockchainClient blockchain.Client) *Service {
	return &Service{
		repo:       repo,
		ipfs:       ipfs,
		bigchain:   bigchain,
		blockchain: blockchainClient,
	}
}

func (s *Service) Upload(ctx context.Context, file multipart.File, fileName string) (models.VideoUploadResponse, error) {
	videoPath, err := writeTempVideo(fileName, file)
	if err != nil {
		return models.VideoUploadResponse{}, err
	}
	defer os.Remove(videoPath)

	audioPath := videoPath + ".mp3"
	defer os.Remove(audioPath)

	if err := audio.ExtractMP3(ctx, videoPath, audioPath); err != nil {
		return models.VideoUploadResponse{}, err
	}

	videoHash, err := hash.SHA256File(videoPath)
	if err != nil {
		return models.VideoUploadResponse{}, fmt.Errorf("hash video: %w", err)
	}
	audioHash, err := hash.SHA256File(audioPath)
	if err != nil {
		return models.VideoUploadResponse{}, fmt.Errorf("hash audio: %w", err)
	}

	videoCID, err := s.ipfs.UploadFile(ctx, videoPath)
	if err != nil {
		return models.VideoUploadResponse{}, fmt.Errorf("upload video ipfs: %w", err)
	}
	audioCID, err := s.ipfs.UploadFile(ctx, audioPath)
	if err != nil {
		return models.VideoUploadResponse{}, fmt.Errorf("upload audio ipfs: %w", err)
	}

	record := models.VideoRecord{
		ProofID:   uuid.NewString(),
		VideoHash: videoHash,
		AudioHash: audioHash,
		VideoCID:  videoCID,
		AudioCID:  audioCID,
		Timestamp: time.Now().UTC(),
		Status:    "PROCESSING",
		FileName:  filepath.Base(fileName),
	}

	receipt, err := s.blockchain.AnchorVideoProof(ctx, record)
	if err != nil {
		return models.VideoUploadResponse{}, fmt.Errorf("anchor video proof: %w", err)
	}
	record.TxHash = receipt.TxHash
	record.Status = "ANCHORED"

	if err := s.repo.SaveVideo(ctx, record); err != nil {
		return models.VideoUploadResponse{}, fmt.Errorf("save video metadata: %w", err)
	}
	if err := s.bigchain.SaveVideoMetadata(ctx, record); err != nil {
		return models.VideoUploadResponse{}, fmt.Errorf("save bigchaindb video metadata: %w", err)
	}

	return models.VideoUploadResponse{
		ProofID:   record.ProofID,
		VideoHash: record.VideoHash,
		AudioHash: record.AudioHash,
		CIDVideo:  record.VideoCID,
		CIDAudio:  record.AudioCID,
		TxHash:    receipt.TxHash,
		Status:    record.Status,
	}, nil
}

func (s *Service) Verify(ctx context.Context, file multipart.File, fileName string) (models.VideoVerifyResponse, error) {
	videoPath, err := writeTempVideo(fileName, file)
	if err != nil {
		return models.VideoVerifyResponse{}, err
	}
	defer os.Remove(videoPath)

	audioPath := videoPath + ".mp3"
	defer os.Remove(audioPath)

	if err := audio.ExtractMP3(ctx, videoPath, audioPath); err != nil {
		return models.VideoVerifyResponse{}, err
	}

	videoHash, err := hash.SHA256File(videoPath)
	if err != nil {
		return models.VideoVerifyResponse{}, fmt.Errorf("hash video: %w", err)
	}
	audioHash, err := hash.SHA256File(audioPath)
	if err != nil {
		return models.VideoVerifyResponse{}, fmt.Errorf("hash audio: %w", err)
	}

	baseRecord, err := s.repo.FindVideoByVideoHash(ctx, videoHash)
	if err != nil {
		return models.VideoVerifyResponse{}, err
	}
	if baseRecord == nil {
		baseRecord, err = s.repo.FindVideoByAudioHash(ctx, audioHash)
		if err != nil {
			return models.VideoVerifyResponse{}, err
		}
	}

	resp := models.VideoVerifyResponse{
		VideoHash: videoHash,
		AudioHash: audioHash,
	}

	if baseRecord == nil {
		resp.Verdict = "TAMPERED"
		resp.Details.Video = "MISMATCH"
		resp.Details.Audio = "MISMATCH"
		return resp, nil
	}

	videoMatch := baseRecord.VideoHash == videoHash
	audioMatch := baseRecord.AudioHash == audioHash

	if videoMatch {
		resp.Details.Video = "MATCH"
	} else {
		resp.Details.Video = "MISMATCH"
	}
	if audioMatch {
		resp.Details.Audio = "MATCH"
	} else {
		resp.Details.Audio = "MISMATCH"
	}

	switch {
	case videoMatch && audioMatch:
		resp.Verdict = "AUTHENTIC"
	case videoMatch && !audioMatch:
		resp.Verdict = "AUDIO_TAMPERED"
	case !videoMatch && audioMatch:
		resp.Verdict = "VIDEO_TAMPERED"
	default:
		resp.Verdict = "TAMPERED"
	}

	return resp, nil
}

func writeTempVideo(fileName string, src multipart.File) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		ext = ".mp4"
	}
	f, err := os.CreateTemp("", "aletheia-video-*"+ext)
	if err != nil {
		return "", fmt.Errorf("create temp video file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, io.LimitReader(src, 2<<30)); err != nil {
		return "", fmt.Errorf("write temp video file: %w", err)
	}
	return f.Name(), nil
}
