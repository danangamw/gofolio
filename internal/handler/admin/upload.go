package admin

import "net/http"

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// TODO: handle image upload (JPG, PNG, WebP), max 5MB, rename with UUID
}
