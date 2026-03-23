package url

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

const defaultMaxImageBytes int64 = 20 << 20

type Service struct {
	httpClient *http.Client
	maxBytes   int64
}

func NewService(timeout time.Duration, maxBytes int64) *Service {
	if timeout <= 0 {
		timeout = 12 * time.Second
	}
	if maxBytes <= 0 {
		maxBytes = defaultMaxImageBytes
	}
	return &Service{
		httpClient: &http.Client{Timeout: timeout},
		maxBytes:   maxBytes,
	}
}

func (s *Service) FetchImageBytes(ctx context.Context, imageURL string) ([]byte, error) {
	parsed, err := neturl.Parse(strings.TrimSpace(imageURL))
	if err != nil {
		return nil, fmt.Errorf("invalid imageUrl: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New("imageUrl must use http or https")
	}
	if parsed.Host == "" {
		return nil, errors.New("imageUrl host is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Aletheia-Verify/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("image fetch failed with status %d", resp.StatusCode)
	}

	contentType := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))
	if contentType != "" && !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "application/octet-stream") {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	limited := io.LimitReader(resp.Body, s.maxBytes+1)
	payload, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read image body: %w", err)
	}
	if int64(len(payload)) > s.maxBytes {
		return nil, fmt.Errorf("image exceeds max size of %d bytes", s.maxBytes)
	}
	if len(payload) == 0 {
		return nil, errors.New("image payload is empty")
	}

	return payload, nil
}
