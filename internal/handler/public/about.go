package public

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-cms/internal/dto"
	"go-cms/internal/repository"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any)
}

type AboutHandler struct {
	tmpl Renderer
	repo *repository.SysConfigRepository
}

func NewAboutHandler(tmpl Renderer, repo *repository.SysConfigRepository) *AboutHandler {
	return &AboutHandler{tmpl: tmpl, repo: repo}
}

func (h *AboutHandler) Index(w http.ResponseWriter, r *http.Request) {
	data := h.defaultAboutData()

	configs, err := h.repo.FindByKeys(r.Context(), []string{
		"about_page_title",
		"about_page_subtitle",
		"about_name",
		"about_role",
		"about_avatar",
		"about_bio_1",
		"about_bio_2",
		"about_skills",
		"about_experiences",
		"about_github_url",
		"about_linkedin_url",
	})
	if err == nil {
		data["Title"] = getString(configs, "about_page_title", "About Me")
		data["AboutTitle"] = data["Title"]
		data["AboutSubtitle"] = getString(configs, "about_page_subtitle", "")
		data["AboutName"] = getString(configs, "about_name", "Danang")
		data["AboutRole"] = getString(configs, "about_role", "Backend Developer")
		data["AboutAvatar"] = getString(configs, "about_avatar", "D")
		data["BioParagraphs"] = []string{
			getString(configs, "about_bio_1", ""),
			getString(configs, "about_bio_2", ""),
		}
		data["AboutSkills"] = parseStringSlice(getString(configs, "about_skills", ""))
		data["AboutExperiences"] = parseExperiences(getString(configs, "about_experiences", ""))
		data["GitHubURL"] = getString(configs, "about_github_url", "https://github.com")
		data["LinkedInURL"] = getString(configs, "about_linkedin_url", "https://linkedin.com")
	}

	h.tmpl.Render(w, "about", data)
}

type aboutExperience struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Period      string `json:"period"`
	Description string `json:"description"`
}

func (h *AboutHandler) defaultAboutData() map[string]any {
	return map[string]any{
		"Title":         "About Me — danangamw",
		"ActiveMenu":    "about",
		"AboutTitle":    "About Me",
		"AboutSubtitle": "Get to know me better, my background, technical expertise, and what I enjoy doing.",
		"AboutName":     "Danang",
		"AboutRole":     "Backend Developer",
		"AboutAvatar":   "D",
		"BioParagraphs": []string{
			"Hello! My name is Danang. I am a Software Engineer specializing in backend development using Go. I love building systems that are modular, clean, well-documented, and efficient.",
			"Besides writing code, I enjoy exploring microservices architectures, database optimization, reliable automated testing, and deploying container-based systems (Docker/Kubernetes).",
		},
		"AboutSkills": []string{
			"Go (Golang)",
			"PostgreSQL",
			"Redis",
			"Docker",
			"Kubernetes",
			"gRPC / Protobuf",
			"GraphQL",
			"Git",
			"HTML / CSS",
		},
		"AboutExperiences": []aboutExperience{
			{
				Title:       "Senior Backend Engineer",
				Company:     "TechCorp Indonesia",
				Period:      "2024 - Present",
				Description: "Designed and migrated monolithic architectures to Go-based microservices, reducing API latency by 40% and handling high-throughput workloads.",
			},
			{
				Title:       "Software Engineer (Backend)",
				Company:     "Startup Hub",
				Period:      "2022 - 2024",
				Description: "Built RESTful APIs with Gin & GORM, integrated third-party payment gateways, and optimized PostgreSQL queries.",
			},
		},
		"GitHubURL":   "https://github.com",
		"LinkedInURL": "https://linkedin.com",
	}
}

func getString(configs map[string]dto.SysConfigResponse, key, fallback string) string {
	if cfg, ok := configs[key]; ok && cfg.Value != "" {
		return cfg.Value
	}
	return fallback
}

func parseStringSlice(raw string) []string {
	if raw == "" {
		return nil
	}
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err == nil {
		return values
	}
	// fallback: comma-separated values
	parts := []string{}
	for _, part := range strings.Split(raw, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func parseExperiences(raw string) []aboutExperience {
	if raw == "" {
		return nil
	}
	var experiences []aboutExperience
	if err := json.Unmarshal([]byte(raw), &experiences); err == nil {
		return experiences
	}
	return nil
}
