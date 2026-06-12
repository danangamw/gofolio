package admin

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-cms/internal/service"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

type uploadResponseData struct {
	FilePath string `json:"filePath"`
}

type uploadResponse struct {
	Data    uploadResponseData `json:"data"`
	URL     string             `json:"url"`
	Status  int                `json:"status"`
	Message string             `json:"message,omitempty"`
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Parse Multipart Form - Max 5MB
	const maxUploadSize = 5 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		h.respondWithError(w, "File size exceeds limit of 5MB", http.StatusBadRequest)
		return
	}

	// Retrieve file from multipart form (support both 'image' and 'file' fields)
	file, header, err := r.FormFile("image")
	if err != nil {
		file, header, err = r.FormFile("file")
		if err != nil {
			h.respondWithError(w, "No file provided under field 'image' or 'file'", http.StatusBadRequest)
			return
		}
	}
	defer file.Close()

	// Check MIME type / extension
	contentType := header.Header.Get("Content-Type")
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}
	if !validTypes[strings.ToLower(contentType)] {
		h.respondWithError(w, "Invalid file format. Only JPEG, PNG, WebP, and GIF are allowed.", http.StatusBadRequest)
		return
	}

	// Perform Upload
	result, err := h.uploadService.Upload(r.Context(), file, "images", header.Filename)
	if err != nil {
		h.respondWithError(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(uploadResponse{
		Data: uploadResponseData{
			FilePath: result.URL,
		},
		URL:    result.URL,
		Status: 1,
	})
}

func (h *UploadHandler) respondWithError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(uploadResponse{
		Status:  0,
		Message: msg,
	})
}
