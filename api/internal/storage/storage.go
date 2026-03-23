package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"aletheia-api/internal/models"
	shell "github.com/ipfs/go-ipfs-api"
)

type IPFSClient interface {
	Upload(ctx context.Context, fileName string, payload []byte) (string, error)
	UploadFile(ctx context.Context, filePath string) (string, error)
}

type MetadataRepository interface {
	Save(ctx context.Context, record models.ProofRecord) error
	FindBySHA256(ctx context.Context, hash string) (*models.ProofRecord, error)
	List(ctx context.Context) ([]models.ProofRecord, error)
	SaveVideo(ctx context.Context, record models.VideoRecord) error
	FindVideoByVideoHash(ctx context.Context, videoHash string) (*models.VideoRecord, error)
	FindVideoByAudioHash(ctx context.Context, audioHash string) (*models.VideoRecord, error)
}

type BigchainClient interface {
	SaveMetadata(ctx context.Context, record models.ProofRecord) error
	SaveVideoMetadata(ctx context.Context, record models.VideoRecord) error
}

type IPFSService struct {
	shell *shell.Shell
}

func NewIPFSService(endpoint string) *IPFSService {
	if endpoint == "" {
		return &IPFSService{}
	}
	return &IPFSService{shell: shell.NewShell(endpoint)}
}

func (s *IPFSService) Upload(ctx context.Context, fileName string, payload []byte) (string, error) {
	if len(payload) == 0 {
		return "", errors.New("empty payload")
	}

	if s.shell == nil {
		// Deterministic fallback CID for local development.
		sum := sha256.Sum256(payload)
		return "local-" + hex.EncodeToString(sum[:])[:32], nil
	}

	reader := bytes.NewReader(payload)
	cid, err := s.shell.Add(reader)
	if err != nil {
		return "", fmt.Errorf("ipfs upload failed: %w", err)
	}
	return cid, nil
}

func (s *IPFSService) UploadFile(ctx context.Context, filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	if s.shell == nil {
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return "", fmt.Errorf("hash file for local cid: %w", err)
		}
		sum := h.Sum(nil)
		return "local-" + hex.EncodeToString(sum)[:32], nil
	}

	cid, err := s.shell.Add(f)
	if err != nil {
		return "", fmt.Errorf("ipfs upload file failed: %w", err)
	}
	return cid, nil
}

type InMemoryRepository struct {
	mu           sync.RWMutex
	byID         map[string]models.ProofRecord
	byHash       map[string]string
	videoByID    map[string]models.VideoRecord
	videoByHash  map[string]string
	audioByHash  map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		byID:        make(map[string]models.ProofRecord),
		byHash:      make(map[string]string),
		videoByID:   make(map[string]models.VideoRecord),
		videoByHash: make(map[string]string),
		audioByHash: make(map[string]string),
	}
}

func (r *InMemoryRepository) Save(_ context.Context, record models.ProofRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[record.ProofID] = record
	r.byHash[record.SHA256] = record.ProofID
	return nil
}

func (r *InMemoryRepository) FindBySHA256(_ context.Context, hash string) (*models.ProofRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byHash[hash]
	if !ok {
		return nil, nil
	}
	record, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	copy := record
	return &copy, nil
}

func (r *InMemoryRepository) List(_ context.Context) ([]models.ProofRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.ProofRecord, 0, len(r.byID))
	for _, v := range r.byID {
		out = append(out, v)
	}
	return out, nil
}

func (r *InMemoryRepository) SaveVideo(_ context.Context, record models.VideoRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.videoByID[record.ProofID] = record
	r.videoByHash[record.VideoHash] = record.ProofID
	r.audioByHash[record.AudioHash] = record.ProofID
	return nil
}

func (r *InMemoryRepository) FindVideoByVideoHash(_ context.Context, videoHash string) (*models.VideoRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.videoByHash[videoHash]
	if !ok {
		return nil, nil
	}
	record, ok := r.videoByID[id]
	if !ok {
		return nil, nil
	}
	copy := record
	return &copy, nil
}

func (r *InMemoryRepository) FindVideoByAudioHash(_ context.Context, audioHash string) (*models.VideoRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.audioByHash[audioHash]
	if !ok {
		return nil, nil
	}
	record, ok := r.videoByID[id]
	if !ok {
		return nil, nil
	}
	copy := record
	return &copy, nil
}

type BigchainHTTPClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewBigchainHTTPClient(baseURL, apiKey string) *BigchainHTTPClient {
	return &BigchainHTTPClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *BigchainHTTPClient) SaveMetadata(ctx context.Context, record models.ProofRecord) error {
	if c.baseURL == "" {
		return nil
	}

	body, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal bigchaindb payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/metadata", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create bigchaindb request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("send bigchaindb request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bigchaindb save failed with status %d", resp.StatusCode)
	}

	return nil
}

func (c *BigchainHTTPClient) SaveVideoMetadata(ctx context.Context, record models.VideoRecord) error {
	if c.baseURL == "" {
		return nil
	}

	body, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal bigchaindb video payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/metadata/video", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create bigchaindb video request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("send bigchaindb video request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bigchaindb save video failed with status %d", resp.StatusCode)
	}

	return nil
}
