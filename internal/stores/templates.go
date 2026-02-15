// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package stores

import (
	"blackoutbox/internal/models"
	"database/sql"
	"time"
)

type TemplateStoreInterface interface {
	Add(model models.Template) error
	Get() ([]models.Template, error)
	Update(model models.Template) error
	GetById(id string) (*models.Template, error)
	GetByFileId(id string) (*models.Template, error)
	GetBySystemId(id string) ([]models.Template, error)
}

type TemplateStore struct {
	Db *sql.DB
}

func (s *TemplateStore) Add(model models.Template) error {
	now := time.Now().Unix()

	_, err := s.Db.Exec(`
		INSERT INTO templates (system_id, file_id, file_path, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, model.SystemId, model.FileId, model.FilePath, model.Description, now, now)
	if err != nil {
		return err
	}

	return nil
}

func (s *TemplateStore) Get() ([]models.Template, error) {
	query, err := s.Db.Query(`
		SELECT id, system_id, file_id, file_path, description, created_at, deleted_at
		FROM templates
	`)
	if err != nil {
		return nil, err
	}

	var documents []models.Template

	for query.Next() {
		var document models.Template

		err := query.Scan(
			&document.Id,
			&document.SystemId,
			&document.FileId,
			&document.FilePath,
			&document.CreatedAt,
			&document.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	return documents, nil
}

func (s *TemplateStore) Update(model models.Template) error {
	return nil
}

func (s *TemplateStore) GetById(id string) (*models.Template, error) {
	row := s.Db.QueryRow(`
		SELECT id, system_id, file_id, file_path, description, deleted_at
		FROM templates
		WHERE id = ?
	`, id)

	var document models.Template

	err := row.Scan(
		&document.Id,
		&document.SystemId,
		&document.FileId,
		&document.FilePath,
		&document.Description,
		&document.CreatedAt,
		&document.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (s *TemplateStore) GetByFileId(id string) (*models.Template, error) {
	row := s.Db.QueryRow(`
		SELECT id, system_id, file_id, file_path, description, created_at, deleted_at
		FROM templates
		WHERE file_id = ?
	`, id)

	var document models.Template

	err := row.Scan(
		&document.Id,
		&document.SystemId,
		&document.FileId,
		&document.FilePath,
		&document.Description,
		&document.CreatedAt,
		&document.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (s *TemplateStore) GetBySystemId(id string) ([]models.Template, error) {
	query, err := s.Db.Query(`
		SELECT id, system_id, file_id, file_path, description, created_at, deleted_at
		FROM templates
		WHERE system_id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	var documents []models.Template

	for query.Next() {
		var document models.Template

		err := query.Scan(
			&document.Id,
			&document.SystemId,
			&document.FileId,
			&document.FilePath,
			&document.Description,
			&document.CreatedAt,
			&document.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	return documents, nil
}
