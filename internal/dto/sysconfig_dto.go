package dto

import "time"

type CreateSysConfigRequest struct {
	Key         string `json:"key" validate:"required,max=255"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	Value       string `json:"value" validate:"required"`
	ValueType   string `json:"value_type" validate:"required,oneof=text textarea url json boolean image"`
	GroupName   string `json:"group_name" validate:"required,max=100"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateSysConfigRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	Value       string `json:"value" validate:"required"`
	ValueType   string `json:"value_type" validate:"required,oneof=text textarea url json boolean image"`
	GroupName   string `json:"group_name" validate:"required,max=100"`
	SortOrder   int    `json:"sort_order"`
}

type SysConfigResponse struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	ValueType   string    `json:"value_type"`
	GroupName   string    `json:"group_name"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c SysConfigResponse) DisplayValue() string {
	switch c.ValueType {
	case "textarea", "json":
		if len(c.Value) > 120 {
			return c.Value[:120] + "..."
		}
		return c.Value
	default:
		if len(c.Value) > 80 {
			return c.Value[:80] + "..."
		}
		return c.Value
	}
}
