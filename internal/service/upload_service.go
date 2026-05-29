package service

// UploadService handles file upload logic (local or S3-compatible storage).
type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

// TODO: Upload (validate MIME type, rename with UUID, save to local or S3)
