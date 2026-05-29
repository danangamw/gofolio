package admin

import "net/http"

type AdminBlogHandler struct{}

func NewAdminBlogHandler() *AdminBlogHandler {
	return &AdminBlogHandler{}
}

func (h *AdminBlogHandler) List(w http.ResponseWriter, r *http.Request)   {}
func (h *AdminBlogHandler) New(w http.ResponseWriter, r *http.Request)    {}
func (h *AdminBlogHandler) Create(w http.ResponseWriter, r *http.Request) {}
func (h *AdminBlogHandler) Edit(w http.ResponseWriter, r *http.Request)   {}
func (h *AdminBlogHandler) Update(w http.ResponseWriter, r *http.Request) {}
func (h *AdminBlogHandler) Delete(w http.ResponseWriter, r *http.Request) {}
