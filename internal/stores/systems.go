// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package stores

import (
	"blackoutbox/internal/models"
	"blackoutbox/internal/storage"
	"database/sql"
	"os"
	"path/filepath"
	"time"
)

type SystemStoreInterface interface {
	Sync(systemId int64, documents []models.Document) error
	AddSystem(system models.System) error
	GetSystems() ([]models.System, error)
	GetSystemById(id int64) (*models.System, error)
	GetSystemByReference(reference string) (*models.System, error)
	UpdateSystem(system models.System) error
	DeleteSystem(id int64) error
}

type SystemStore struct {
	Db *sql.DB
}

// Sync replaces all documents for a system and refreshes its filesystem folder.
func (s *SystemStore) Sync(systemId int64, documents []models.Document) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Get the system reference for filesystem operations
	var systemRef string
	err = s.Db.QueryRow(`SELECT reference FROM systems WHERE id = ?`, systemId).Scan(&systemRef)
	if err != nil {
		return err
	}

	// 2. Remove existing documents for the system
	_, err = tx.Exec(`
		DELETE FROM documents
		WHERE system_id = ?
	`, systemId)
	if err != nil {
		return err
	}

	// 3. Insert new document metadata
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

	now := time.Now().Unix()

	for _, doc := range documents {
		_, err := stmt.Exec(
			systemId,
			doc.FileReference,
			doc.FilePath,
			doc.PrintAt,
			doc.LastPrintedAt,
			now,
		)
		if err != nil {
			return err
		}
	}

	// 4. Commit database changes
	if err := tx.Commit(); err != nil {
		return err
	}

	// 5. Sync filesystem
	systemDir := filepath.Join(storage.DocumentsRoot, systemRef)

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

func (s *SystemStore) AddSystem(system models.System) error {
	_, err := s.Db.Exec(`
		INSERT INTO systems (reference, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, system.Reference, system.Name, system.Description, system.CreatedAt, system.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *SystemStore) GetSystems() ([]models.System, error) {
	query, err := s.Db.Query(`
		SELECT id, reference, name, description, created_at, updated_at, deleted_at
		FROM systems
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var systems []models.System

	for query.Next() {
		var system models.System

		err := query.Scan(
			&system.Id,
			&system.Reference,
			&system.Name,
			&system.Description,
			&system.CreatedAt,
			&system.UpdatedAt,
			system.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		systems = append(systems, system)
	}

	return systems, nil
}

func (s *SystemStore) GetSystemById(id int64) (*models.System, error) {
	row := s.Db.QueryRow(`
		SELECT id, reference, name, description, created_at, updated_at, deleted_at
		FROM systems
		WHERE id = ? AND deleted_at IS NULL
	`, id)

	var system models.System

	err := row.Scan(
		&system.Id,
		&system.Reference,
		&system.Name,
		&system.Description,
		&system.CreatedAt,
		&system.UpdatedAt,
		system.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &system, nil
}

func (s *SystemStore) GetSystemByReference(reference string) (*models.System, error) {
	row := s.Db.QueryRow(`
		SELECT id, reference, name, description, created_at, updated_at, deleted_at
		FROM systems
		WHERE reference = ? AND deleted_at IS NULL
	`, reference)

	var system models.System

	err := row.Scan(
		&system.Id,
		&system.Reference,
		&system.Name,
		&system.Description,
		&system.CreatedAt,
		&system.UpdatedAt,
		system.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &system, nil
}

func (s *SystemStore) UpdateSystem(system models.System) error {
	_, err := s.Db.Exec(`
		UPDATE systems
		SET reference = ?, name = ?, description = ?, updated_at = ?
		WHERE id = ?
	`, system.Reference, system.Name, system.Description, time.Now().Unix(), system.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SystemStore) DeleteSystem(id int64) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the system reference for filesystem operations
	var systemRef string
	err = s.Db.QueryRow(`SELECT reference FROM systems WHERE id = ?`, id).Scan(&systemRef)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM documents
		WHERE system_id = ?
	`, id)
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	_, err = s.Db.Exec(`
		UPDATE systems
		SET deleted_at = ?
		WHERE id = ?
	`, now, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	systemDir := filepath.Join(storage.DocumentsRoot, systemRef)
	return os.RemoveAll(systemDir)
}
