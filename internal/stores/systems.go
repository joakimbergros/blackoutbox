package stores

import (
	"blackoutbox/internal/models"
	"database/sql"
	"os"
	"path/filepath"
	"time"
)

type SystemStoreInterface interface {
	Sync(systemId string, documents []models.Document) error
	DeleteSystem(systemId string) error
}

type SystemStore struct {
	Db        *sql.DB
	FilesRoot string // root directory where system folders live
}

// Sync replaces all documents for a system and refreshes its filesystem folder.
func (s *SystemStore) Sync(systemId string, documents []models.Document) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Remove existing documents for the system
	_, err = tx.Exec(`
		DELETE FROM documents
		WHERE system_id = ?
	`, systemId)
	if err != nil {
		return err
	}

	// 2. Insert new document metadata
	stmt, err := tx.Prepare(`
		INSERT INTO documents (
			system_id,
			file_id,
			file_path,
			print_at,
			last_printed_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()

	for _, doc := range documents {
		_, err := stmt.Exec(
			systemId,
			doc.FileId,
			doc.FilePath,
			doc.PrintAt,
			doc.LastPrintedAt,
			&now,
		)
		if err != nil {
			return err
		}
	}

	// 3. Commit database changes
	if err := tx.Commit(); err != nil {
		return err
	}

	// 4. Sync filesystem
	systemDir := filepath.Join(s.FilesRoot, systemId)

	// Remove old system directory completely
	if err := os.RemoveAll(systemDir); err != nil {
		return err
	}

	// Recreate system directory
	if err := os.MkdirAll(systemDir, 0755); err != nil {
		return err
	}

	return nil
}

// DeleteSystem removes all documents and filesystem data for a system.
func (s *SystemStore) DeleteSystem(systemId string) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM documents
		WHERE system_id = ?
	`, systemId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	systemDir := filepath.Join(s.FilesRoot, systemId)
	return os.RemoveAll(systemDir)
}
