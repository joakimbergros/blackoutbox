// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package stores

import (
	"blackoutbox/internal/models"
	"database/sql"
	"encoding/json"
	"time"
)

type DocumentStoreInterface interface {
	Add(model models.Document) error
	Get() ([]models.Document, error)
	Update(model models.Document) error
	GetById(id int64) (*models.Document, error)
	GetByFileId(id int64) (*models.Document, error)
	GetBySystemId(id int64) ([]models.Document, error)
}

type DocumentStore struct {
	Db *sql.DB
}

func (s *DocumentStore) Add(model models.Document) error {
	tagsJSON, err := json.Marshal(model.Tags)
	if err != nil {
		return err
	}

	updatedAt := time.Now().Unix()

	_, err = s.Db.Exec(`
		INSERT INTO documents (system_id, file_id, file_path, print_at, last_printed_at, tags, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, model.SystemId, model.FileReference, model.FilePath, model.PrintAt, model.LastPrintedAt, string(tagsJSON), updatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *DocumentStore) Get() ([]models.Document, error) {
	query, err := s.Db.Query(`
		SELECT id, system_id, file_id, file_path, print_at, last_printed_at, tags, updated_at, deleted_at
		FROM documents
	`)
	if err != nil {
		return nil, err
	}

	var documents []models.Document

	for query.Next() {
		var document models.Document
		var tagsJSON string

		err := query.Scan(
			&document.Id,
			&document.SystemId,
			&document.FileReference,
			&document.FilePath,
			&document.PrintAt,
			&document.LastPrintedAt,
			&tagsJSON,
			&document.UpdatedAt,
			&document.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &document.Tags); err != nil {
				return nil, err
			}
		}

		documents = append(documents, document)
	}

	return documents, nil
}

func (s *DocumentStore) Update(model models.Document) error {
	return nil
}

func (s *DocumentStore) GetById(id int64) (*models.Document, error) {
	row := s.Db.QueryRow(`
		SELECT id, system_id, file_id, file_path, print_at, last_printed_at, tags, updated_at, deleted_at
		FROM documents
		WHERE id = ?
	`, id)

	var document models.Document
	var tagsJSON string

	err := row.Scan(
		&document.Id,
		&document.SystemId,
		&document.FileReference,
		&document.FilePath,
		&document.PrintAt,
		&document.LastPrintedAt,
		&tagsJSON,
		&document.UpdatedAt,
		&document.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &document.Tags); err != nil {
			return nil, err
		}
	}

	return &document, nil
}

func (s *DocumentStore) GetByFileId(id int64) (*models.Document, error) {
	row := s.Db.QueryRow(`
		SELECT id, system_id, file_id, file_path, print_at, last_printed_at, tags, updated_at, deleted_at
		FROM documents
		WHERE file_id = ?
	`, id)

	var document models.Document
	var tagsJSON string

	err := row.Scan(
		&document.Id,
		&document.SystemId,
		&document.FileReference,
		&document.FilePath,
		&document.PrintAt,
		&document.LastPrintedAt,
		&tagsJSON,
		&document.UpdatedAt,
		&document.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &document.Tags); err != nil {
			return nil, err
		}
	}

	return &document, nil
}

func (s *DocumentStore) GetBySystemId(id int64) ([]models.Document, error) {
	query, err := s.Db.Query(`
		SELECT id, system_id, file_id, file_path, print_at, last_printed_at, tags, updated_at, deleted_at
		FROM documents
		WHERE system_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var documents []models.Document

	for query.Next() {
		var document models.Document
		var tagsJSON string

		err := query.Scan(
			&document.Id,
			&document.SystemId,
			&document.FileReference,
			&document.FilePath,
			&document.PrintAt,
			&document.LastPrintedAt,
			&tagsJSON,
			&document.UpdatedAt,
			&document.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &document.Tags); err != nil {
				return nil, err
			}
		}

		documents = append(documents, document)
	}

	return documents, nil
}
