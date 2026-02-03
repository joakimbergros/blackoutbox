package models

import "time"

type Document struct {
	Id         int        `json:"id"`
	ExternalId string     `json:"ext_id"`
	IsSystem   bool       `json:"is_system"`
	FilePath   string     `json:"file_pat"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
