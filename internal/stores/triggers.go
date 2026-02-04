package stores

import (
	"blackoutbox/internal/models"
	"database/sql"
	"time"
)

type TriggerStoreInterface interface {
	Add(model models.Trigger) error
	Get() ([]models.Trigger, error)
	GetById(id string) (*models.Trigger, error)
	GetBySystemId(id string) (*models.Trigger, error)
	Update(model models.Trigger) error
	UpdateStatus(id string, status string) error
	IncrementRetryCount(id string) error
	ResetRetryCount(id string) error
	Delete(id string) error
}

type TriggerStore struct {
	Db *sql.DB
}

func (s *TriggerStore) Add(model models.Trigger) error {
	_, err := s.Db.Exec(`
		INSERT INTO triggers (system_id, url, last_failed_at, buffer_seconds, status, last_checked_at, retry_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, model.SystemId, model.Url, model.LastFailedAt, model.BufferSeconds, model.Status, model.LastCheckedAt, model.RetryCount, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (s *TriggerStore) Get() ([]models.Trigger, error) {
	query, err := s.Db.Query(`
		SELECT id, system_id, url, last_failed_at, buffer_seconds, status, last_checked_at, retry_count, created_at, updated_at
		FROM triggers
	`)
	if err != nil {
		return nil, err
	}

	var triggers []models.Trigger

	for query.Next() {
		var trigger models.Trigger

		err := query.Scan(
			&trigger.Id,
			&trigger.SystemId,
			&trigger.Url,
			&trigger.LastFailedAt,
			&trigger.BufferSeconds,
			&trigger.Status,
			&trigger.LastCheckedAt,
			&trigger.RetryCount,
			&trigger.CreatedAt,
			&trigger.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		triggers = append(triggers, trigger)
	}

	return triggers, nil
}

func (s *TriggerStore) GetById(id string) (*models.Trigger, error) {
	row := s.Db.QueryRow(`
		SELECT id, system_id, url, last_failed_at, buffer_seconds, status, last_checked_at, retry_count, created_at, updated_at
		FROM triggers
		WHERE id = ?
	`, id)

	var trigger models.Trigger

	err := row.Scan(
		&trigger.Id,
		&trigger.SystemId,
		&trigger.Url,
		&trigger.LastFailedAt,
		&trigger.BufferSeconds,
		&trigger.Status,
		&trigger.LastCheckedAt,
		&trigger.RetryCount,
		&trigger.CreatedAt,
		&trigger.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &trigger, nil
}

func (s *TriggerStore) GetBySystemId(id string) (*models.Trigger, error) {
	row := s.Db.QueryRow(`
		SELECT id, system_id, url, last_failed_at, buffer_seconds, status, last_checked_at, retry_count, created_at, updated_at
		FROM triggers
		WHERE system_id = ?
	`, id)

	var trigger models.Trigger

	err := row.Scan(
		&trigger.Id,
		&trigger.SystemId,
		&trigger.Url,
		&trigger.LastFailedAt,
		&trigger.BufferSeconds,
		&trigger.Status,
		&trigger.LastCheckedAt,
		&trigger.RetryCount,
		&trigger.CreatedAt,
		&trigger.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &trigger, nil
}

func (s *TriggerStore) Update(model models.Trigger) error {
	_, err := s.Db.Exec(`
		UPDATE triggers
		SET system_id = ?, url = ?, last_failed_at = ?, buffer_seconds = ?, status = ?, last_checked_at = ?, retry_count = ?, updated_at = ?
		WHERE id = ?
	`, model.SystemId, model.Url, model.LastFailedAt, model.BufferSeconds, model.Status, model.LastCheckedAt, model.RetryCount, time.Now(), model.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TriggerStore) UpdateStatus(id string, status string) error {
	_, err := s.Db.Exec(`
		UPDATE triggers
		SET status = ?, updated_at = ?
		WHERE id = ?
	`, status, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TriggerStore) IncrementRetryCount(id string) error {
	_, err := s.Db.Exec(`
		UPDATE triggers
		SET retry_count = retry_count + 1, updated_at = ?
		WHERE id = ?
	`, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TriggerStore) ResetRetryCount(id string) error {
	_, err := s.Db.Exec(`
		UPDATE triggers
		SET retry_count = 0, last_failed_at = NULL, status = 'ok', updated_at = ?
		WHERE id = ?
	`, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TriggerStore) Delete(id string) error {
	_, err := s.Db.Exec("DELETE FROM triggers WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
