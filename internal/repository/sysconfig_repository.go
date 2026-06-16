package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-cms/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SysConfigRepository struct {
	db *gorm.DB
}

func NewSysConfigRepository(db *gorm.DB) *SysConfigRepository {
	return &SysConfigRepository{db: db}
}

func (r *SysConfigRepository) FindAll(ctx context.Context) ([]dto.SysConfigResponse, error) {
	var configs []dto.SysConfigResponse
	query := `
		SELECT id, key, name, description, value, value_type, group_name, sort_order, created_at, updated_at
		FROM sys_configs
		ORDER BY group_name ASC, sort_order ASC, created_at DESC
	`
	if err := r.db.WithContext(ctx).Raw(query).Scan(&configs).Error; err != nil {
		return nil, fmt.Errorf("sysconfig repo: find all: %w", err)
	}
	return configs, nil
}

func (r *SysConfigRepository) FindByKey(ctx context.Context, key string) (dto.SysConfigResponse, error) {
	var config dto.SysConfigResponse
	query := `
		SELECT id, key, name, description, value, value_type, group_name, sort_order, created_at, updated_at
		FROM sys_configs
		WHERE key = ?
	`
	if err := r.db.WithContext(ctx).Raw(query, key).Scan(&config).Error; err != nil {
		return dto.SysConfigResponse{}, fmt.Errorf("sysconfig repo: find by key: %w", err)
	}
	if config.ID == "" {
		return dto.SysConfigResponse{}, gorm.ErrRecordNotFound
	}
	return config, nil
}

func (r *SysConfigRepository) FindByGroup(ctx context.Context, group string) ([]dto.SysConfigResponse, error) {
	var configs []dto.SysConfigResponse
	query := `
		SELECT id, key, name, description, value, value_type, group_name, sort_order, created_at, updated_at
		FROM sys_configs
		WHERE group_name = ?
		ORDER BY sort_order ASC, created_at DESC
	`
	if err := r.db.WithContext(ctx).Raw(query, group).Scan(&configs).Error; err != nil {
		return nil, fmt.Errorf("sysconfig repo: find by group: %w", err)
	}
	return configs, nil
}

func (r *SysConfigRepository) FindByKeys(ctx context.Context, keys []string) (map[string]dto.SysConfigResponse, error) {
	if len(keys) == 0 {
		return map[string]dto.SysConfigResponse{}, nil
	}

	var configs []dto.SysConfigResponse
	query := `
		SELECT id, key, name, description, value, value_type, group_name, sort_order, created_at, updated_at
		FROM sys_configs
		WHERE key IN ?
		ORDER BY group_name ASC, sort_order ASC, created_at DESC
	`
	if err := r.db.WithContext(ctx).Raw(query, keys).Scan(&configs).Error; err != nil {
		return nil, fmt.Errorf("sysconfig repo: find by keys: %w", err)
	}

	resp := make(map[string]dto.SysConfigResponse, len(configs))
	for _, cfg := range configs {
		resp[cfg.Key] = cfg
	}
	return resp, nil
}

func (r *SysConfigRepository) Create(ctx context.Context, req dto.CreateSysConfigRequest) (dto.SysConfigResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return dto.SysConfigResponse{}, fmt.Errorf("sysconfig repo: create: generate uuid: %w", err)
	}
	now := time.Now()

	valueType := normalizeValueType(req.ValueType)
	groupName := normalizeGroupName(req.GroupName)

	query := `
		INSERT INTO sys_configs (id, key, name, description, value, value_type, group_name, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	if err := r.db.WithContext(ctx).Exec(query, id, req.Key, req.Name, req.Description, req.Value, valueType, groupName, req.SortOrder, now, now).Error; err != nil {
		return dto.SysConfigResponse{}, fmt.Errorf("sysconfig repo: create: %w", err)
	}

	return dto.SysConfigResponse{
		ID:          id.String(),
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Value:       req.Value,
		ValueType:   valueType,
		GroupName:   groupName,
		SortOrder:   req.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *SysConfigRepository) Update(ctx context.Context, id string, req dto.UpdateSysConfigRequest) (dto.SysConfigResponse, error) {
	now := time.Now()
	valueType := normalizeValueType(req.ValueType)
	groupName := normalizeGroupName(req.GroupName)

	query := `
		UPDATE sys_configs
		SET name = ?, description = ?, value = ?, value_type = ?, group_name = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`
	result := r.db.WithContext(ctx).Exec(query, req.Name, req.Description, req.Value, valueType, groupName, req.SortOrder, now, id)
	if result.Error != nil {
		return dto.SysConfigResponse{}, fmt.Errorf("sysconfig repo: update: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return dto.SysConfigResponse{}, gorm.ErrRecordNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *SysConfigRepository) FindByID(ctx context.Context, id string) (dto.SysConfigResponse, error) {
	var config dto.SysConfigResponse
	query := `
		SELECT id, key, name, description, value, value_type, group_name, sort_order, created_at, updated_at
		FROM sys_configs
		WHERE id = ?
	`
	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&config).Error; err != nil {
		return dto.SysConfigResponse{}, fmt.Errorf("sysconfig repo: find by id: %w", err)
	}
	if config.ID == "" {
		return dto.SysConfigResponse{}, gorm.ErrRecordNotFound
	}
	return config, nil
}

func (r *SysConfigRepository) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM sys_configs WHERE key = ?`
	result := r.db.WithContext(ctx).Exec(query, key)
	if result.Error != nil {
		return fmt.Errorf("sysconfig repo: delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SysConfigRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM sys_configs`).Scan(&count).Error; err != nil {
		return 0, fmt.Errorf("sysconfig repo: count: %w", err)
	}
	return count, nil
}

func (r *SysConfigRepository) SeedDefaults(ctx context.Context) error {
	defaults := []dto.CreateSysConfigRequest{
		{Key: "site_name", Name: "Site Name", Description: "Primary brand name shown across the site.", Value: "danang.dev", ValueType: "text", GroupName: "general", SortOrder: 1},
		{Key: "site_tagline", Name: "Site Tagline", Description: "Short supporting tagline for the public site.", Value: "Backend Developer Portfolio", ValueType: "text", GroupName: "general", SortOrder: 2},
		{Key: "about_page_title", Name: "About Page Title", Description: "Main title for the About page.", Value: "About Me", ValueType: "text", GroupName: "about", SortOrder: 1},
		{Key: "about_page_subtitle", Name: "About Page Subtitle", Description: "Supporting subtitle for the About page header.", Value: "Get to know me better, my background, technical expertise, and what I enjoy doing.", ValueType: "textarea", GroupName: "about", SortOrder: 2},
		{Key: "about_name", Name: "About Name", Description: "Displayed name in the About profile card.", Value: "Danang", ValueType: "text", GroupName: "about", SortOrder: 3},
		{Key: "about_role", Name: "About Role", Description: "Job title or role shown below the name.", Value: "Backend Developer", ValueType: "text", GroupName: "about", SortOrder: 4},
		{Key: "about_avatar", Name: "About Avatar", Description: "Single-character avatar or emoji for the profile card.", Value: "D", ValueType: "text", GroupName: "about", SortOrder: 5},
		{Key: "about_bio_1", Name: "About Bio Paragraph 1", Description: "First bio paragraph shown on the About page.", Value: "Hello! My name is Danang. I am a Software Engineer specializing in backend development using Go. I love building systems that are modular, clean, well-documented, and efficient.", ValueType: "textarea", GroupName: "about", SortOrder: 6},
		{Key: "about_bio_2", Name: "About Bio Paragraph 2", Description: "Second bio paragraph shown on the About page.", Value: "Besides writing code, I enjoy exploring microservices architectures, database optimization, reliable automated testing, and deploying container-based systems (Docker/Kubernetes).", ValueType: "textarea", GroupName: "about", SortOrder: 7},
		{Key: "about_skills", Name: "About Skills", Description: "JSON array of skills shown as badges.", Value: `["Go (Golang)","PostgreSQL","Redis","Docker","Kubernetes","gRPC / Protobuf","GraphQL","Git","HTML / CSS"]`, ValueType: "json", GroupName: "about", SortOrder: 8},
		{Key: "about_experiences", Name: "About Experiences", Description: "JSON array of work experience items.", Value: `[{"title":"Senior Backend Engineer","company":"TechCorp Indonesia","period":"2024 - Present","description":"Designed and migrated monolithic architectures to Go-based microservices, reducing API latency by 40% and handling high-throughput workloads."},{"title":"Software Engineer (Backend)","company":"Startup Hub","period":"2022 - 2024","description":"Built RESTful APIs with Gin & GORM, integrated third-party payment gateways, and optimized PostgreSQL queries."}]`, ValueType: "json", GroupName: "about", SortOrder: 9},
		{Key: "about_github_url", Name: "About GitHub URL", Description: "GitHub profile link shown on the About page.", Value: "https://github.com", ValueType: "url", GroupName: "about", SortOrder: 10},
		{Key: "about_linkedin_url", Name: "About LinkedIn URL", Description: "LinkedIn profile link shown on the About page.", Value: "https://linkedin.com", ValueType: "url", GroupName: "about", SortOrder: 11},
	}

	for _, def := range defaults {
		if err := r.seedDefault(ctx, def); err != nil {
			return err
		}
	}
	return nil
}

func (r *SysConfigRepository) seedDefault(ctx context.Context, req dto.CreateSysConfigRequest) error {
	query := `
		INSERT INTO sys_configs (id, key, name, description, value, value_type, group_name, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (key) DO NOTHING
	`
	now := time.Now()
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("sysconfig repo: seed default: generate uuid: %w", err)
	}
	valueType := normalizeValueType(req.ValueType)
	groupName := normalizeGroupName(req.GroupName)

	if err := r.db.WithContext(ctx).Exec(query, id, req.Key, req.Name, req.Description, req.Value, valueType, groupName, req.SortOrder, now, now).Error; err != nil {
		return fmt.Errorf("sysconfig repo: seed default %q: %w", req.Key, err)
	}
	return nil
}

func normalizeValueType(valueType string) string {
	valueType = strings.ToLower(strings.TrimSpace(valueType))
	switch valueType {
	case "textarea", "url", "json", "boolean", "image":
		return valueType
	default:
		return "text"
	}
}

func normalizeGroupName(groupName string) string {
	groupName = strings.ToLower(strings.TrimSpace(groupName))
	if groupName == "" {
		return "general"
	}
	return groupName
}
