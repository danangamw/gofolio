package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go-cms/internal/config"
	"go-cms/pkg/storage"

	"github.com/google/uuid"
)

// UploadService handles file upload logic (local or S3-compatible storage).
type UploadService struct {
	cfg *config.Config
}

func NewUploadService(cfg *config.Config) *UploadService {
	return &UploadService{cfg: cfg}
}

// Upload Result contains the reference key/path and the public-facing URL
type UploadResult struct {
	Key string // Internal reference (e.g. path on disk or S3 key)
	URL string // Publicly accessible URL (subdomain URL, presigned URL, or local route URL)
}

// Upload uploads the file either to S3/MinIO or to the local disk depending on configuration.
func (s *UploadService) Upload(ctx context.Context, file io.Reader, folder string, originalName string) (*UploadResult, error) {
	ext := filepath.Ext(originalName)
	uniqueName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)

	if s.cfg.UploadStorage == "s3" {
		key, err := storage.UploadFile(ctx, file, folder, uniqueName)
		if err != nil {
			return nil, fmt.Errorf("failed to upload to S3: %w", err)
		}

		publicURL := storage.GetPublicURL(key)
		return &UploadResult{
			Key: key,
			URL: publicURL,
		}, nil
	}

	// Local Storage Upload
	destDir := filepath.Join(s.cfg.UploadDir, folder)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	destPath := filepath.Join(destDir, uniqueName)
	outFile, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create local file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return nil, fmt.Errorf("failed to save local file contents: %w", err)
	}

	// Local URL is served under /uploads/ route
	urlPath := fmt.Sprintf("/uploads/%s/%s", folder, uniqueName)
	return &UploadResult{
		Key: filepath.Join(folder, uniqueName),
		URL: urlPath,
	}, nil
}
