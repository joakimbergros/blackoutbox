// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package stores

import (
	"blackoutbox/internal/models"
	"database/sql"
	"time"
)

type PrintJobStoreInterface interface {
	Add(job models.PrintJob) error
	Get() ([]models.PrintJob, error)
	GetById(id string) (*models.PrintJob, error)
	GetByDocumentId(id string) ([]models.PrintJob, error)
	GetStuckJobs(thresholdSeconds int) ([]models.PrintJob, error)
	Update(job models.PrintJob) error
	UpdateStatus(id string, status string) error
}

type PrintJobStore struct {
	Db *sql.DB
}

func (s *PrintJobStore) Add(job models.PrintJob) error {
	_, err := s.Db.Exec(`
		INSERT INTO print_jobs (document_id, cups_job_id, status, submitted_at, completed_at, error_message)
		VALUES (?, ?, ?, ?, ?, ?)
	`, job.DocumentId, job.CupsJobId, job.Status, job.SubmittedAt, job.CompletedAt, job.ErrorMessage)
	if err != nil {
		return err
	}
	return nil
}

func (s *PrintJobStore) Get() ([]models.PrintJob, error) {
	query, err := s.Db.Query(`
		SELECT id, document_id, cups_job_id, status, submitted_at, completed_at, error_message
		FROM print_jobs
	`)
	if err != nil {
		return nil, err
	}

	var jobs []models.PrintJob

	for query.Next() {
		var job models.PrintJob

		err := query.Scan(
			&job.Id,
			&job.DocumentId,
			&job.CupsJobId,
			&job.Status,
			&job.SubmittedAt,
			&job.CompletedAt,
			&job.ErrorMessage,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *PrintJobStore) GetById(id string) (*models.PrintJob, error) {
	row := s.Db.QueryRow(`
		SELECT id, document_id, cups_job_id, status, submitted_at, completed_at, error_message
		FROM print_jobs
		WHERE id = ?
	`, id)

	var job models.PrintJob

	err := row.Scan(
		&job.Id,
		&job.DocumentId,
		&job.CupsJobId,
		&job.Status,
		&job.SubmittedAt,
		&job.CompletedAt,
		&job.ErrorMessage,
	)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *PrintJobStore) GetByDocumentId(id string) ([]models.PrintJob, error) {
	query, err := s.Db.Query(`
		SELECT id, document_id, cups_job_id, status, submitted_at, completed_at, error_message
		FROM print_jobs
		WHERE document_id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	var jobs []models.PrintJob

	for query.Next() {
		var job models.PrintJob

		err := query.Scan(
			&job.Id,
			&job.DocumentId,
			&job.CupsJobId,
			&job.Status,
			&job.SubmittedAt,
			&job.CompletedAt,
			&job.ErrorMessage,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *PrintJobStore) GetStuckJobs(thresholdSeconds int) ([]models.PrintJob, error) {
	threshold := time.Now().Unix() - int64(thresholdSeconds)

	query, err := s.Db.Query(`
		SELECT id, document_id, cups_job_id, status, submitted_at, completed_at, error_message
		FROM print_jobs
		WHERE status IN ('pending', 'printing') AND submitted_at < ?
	`, threshold)
	if err != nil {
		return nil, err
	}

	var jobs []models.PrintJob

	for query.Next() {
		var job models.PrintJob

		err := query.Scan(
			&job.Id,
			&job.DocumentId,
			&job.CupsJobId,
			&job.Status,
			&job.SubmittedAt,
			&job.CompletedAt,
			&job.ErrorMessage,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *PrintJobStore) Update(job models.PrintJob) error {
	_, err := s.Db.Exec(`
		UPDATE print_jobs
		SET cups_job_id = ?, status = ?, submitted_at = ?, completed_at = ?, error_message = ?
		WHERE id = ?
	`, job.CupsJobId, job.Status, job.SubmittedAt, job.CompletedAt, job.ErrorMessage, job.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PrintJobStore) UpdateStatus(id string, status string) error {
	_, err := s.Db.Exec(`
		UPDATE print_jobs
		SET status = ?
		WHERE id = ?
	`, status, id)
	if err != nil {
		return err
	}
	return nil
}
