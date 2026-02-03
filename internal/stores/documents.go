package stores

import (
	"blackoutbox/internal/models"
	"database/sql"
	"time"
)

type DocumentStoreInterface interface {
	Add(model models.Document) error
	Get() ([]models.Document, error)
	Update(model models.Document) error
	GetByFileId(id int) (*models.Document, error)
	GetBySystemId(id int) ([]models.Document, error)
}

type DocumentStore struct {
	Db *sql.DB
}

func (h *DocumentStore) Add(model models.Document) error {
	_, err := h.Db.Exec(`
		INSERT INTO documents (ext_id, is_system, file_path, updated_at)
		VALUES (?, ?, ?, ?)
	`, model.ExternalId, model.IsSystem, model.FilePath, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (h *DocumentStore) Get() ([]models.Document, error) {
	query, err := h.Db.Query(`
		SELECT * FROM documents
	`)
	if err != nil {
		return nil, err
	}

	var documents []models.Document

	for query.Next() {
		var document models.Document

		err := query.Scan(
			&document.Id,
			&document.ExternalId,
			&document.IsSystem,
			&document.FilePath,
			&document.UpdatedAt,
			&document.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	return documents, nil
}

func (s *DocumentStore) Update(model models.Document) error {
	return nil
}

func (s *DocumentStore) GetByFileId(id int) (*models.Document, error) {
	return nil, nil
}

func (s *DocumentStore) GetBySystemId(id int) ([]models.Document, error) {
	return nil, nil
}
