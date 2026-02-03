package stores

import (
	"blackoutbox/internal/models"
	"database/sql"
	"time"
)

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
