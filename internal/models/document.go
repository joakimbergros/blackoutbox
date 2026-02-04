package models

import "time"

type Document struct {
	Id            int        `json:"id"`
	SystemId      string     `json:"system_id"`
	FileId        string     `json:"file_id"`
	FilePath      string     `json:"file_path"`
	PrintAt       *int64     `json:"print_at"`
	LastPrintedAt *int64     `json:"last_printed_at"`
	Tags          []string   `json:"tags"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
