package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go-cms/internal/dto"
	"go-cms/internal/middleware"
	"go-cms/internal/repository"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AdminSysConfigHandler struct {
	tmpl Renderer
	repo *repository.SysConfigRepository
}

func NewAdminSysConfigHandler(tmpl Renderer, repo *repository.SysConfigRepository) *AdminSysConfigHandler {
	return &AdminSysConfigHandler{
		tmpl: tmpl,
		repo: repo,
	}
}

func (h *AdminSysConfigHandler) List(w http.ResponseWriter, r *http.Request) {
	configs, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":      "Manage Site Settings",
		"ActiveMenu": "sysconfig",
		"Configs":    configs,
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "sysconfig_list", data)
}

func (h *AdminSysConfigHandler) New(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":      "New Site Setting",
		"ActiveMenu": "sysconfig",
		"IsEdit":     false,
		"Config": dto.SysConfigResponse{
			ValueType: "text",
			GroupName: "general",
		},
		"CSRFToken": middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "sysconfig_form", data)
}

func (h *AdminSysConfigHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	req := dto.CreateSysConfigRequest{
		Key:         strings.TrimSpace(r.FormValue("key")),
		Name:        strings.TrimSpace(r.FormValue("name")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Value:       r.FormValue("value"),
		ValueType:   r.FormValue("value_type"),
		GroupName:   r.FormValue("group_name"),
	}

	if order := strings.TrimSpace(r.FormValue("sort_order")); order != "" {
		if parsed, err := strconv.Atoi(order); err == nil {
			req.SortOrder = parsed
		}
	}

	if _, err := h.repo.Create(r.Context(), req); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/sysconfigs", http.StatusSeeOther)
}

func (h *AdminSysConfigHandler) Edit(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	config, err := h.repo.FindByKey(r.Context(), key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":      "Edit Site Setting",
		"ActiveMenu": "sysconfig",
		"IsEdit":     true,
		"Config":     config,
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "sysconfig_form", data)
}

func (h *AdminSysConfigHandler) Update(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	config, err := h.repo.FindByKey(r.Context(), key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	req := dto.UpdateSysConfigRequest{
		Name:        strings.TrimSpace(r.FormValue("name")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Value:       r.FormValue("value"),
		ValueType:   r.FormValue("value_type"),
		GroupName:   r.FormValue("group_name"),
	}
	if order := strings.TrimSpace(r.FormValue("sort_order")); order != "" {
		if parsed, err := strconv.Atoi(order); err == nil {
			req.SortOrder = parsed
		}
	}

	if _, err := h.repo.Update(r.Context(), config.ID, req); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/sysconfigs", http.StatusSeeOther)
}

func (h *AdminSysConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if err := h.repo.Delete(r.Context(), key); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/sysconfigs", http.StatusSeeOther)
}
