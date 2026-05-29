package admin

import "net/http"

type AdminPortfolioHandler struct{}

func NewAdminPortfolioHandler() *AdminPortfolioHandler {
	return &AdminPortfolioHandler{}
}

func (h *AdminPortfolioHandler) List(w http.ResponseWriter, r *http.Request)   {}
func (h *AdminPortfolioHandler) New(w http.ResponseWriter, r *http.Request)    {}
func (h *AdminPortfolioHandler) Create(w http.ResponseWriter, r *http.Request) {}
func (h *AdminPortfolioHandler) Edit(w http.ResponseWriter, r *http.Request)   {}
func (h *AdminPortfolioHandler) Update(w http.ResponseWriter, r *http.Request) {}
func (h *AdminPortfolioHandler) Delete(w http.ResponseWriter, r *http.Request) {}
