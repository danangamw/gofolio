package public

import "net/http"

type BlogHandler struct{}

func NewBlogHandler() *BlogHandler {
	return &BlogHandler{}
}

func (h *BlogHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: render blog_list.html template with pagination
}

func (h *BlogHandler) Detail(w http.ResponseWriter, r *http.Request) {
	// TODO: render blog_detail.html template with markdown-rendered content
}
