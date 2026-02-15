// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package models

type PrintJob struct {
	Id           int64   `json:"id"`
	DocumentId   int64   `json:"document_id"`
	CupsJobId    *string `json:"cups_job_id"`
	Status       string  `json:"status"` // pending, printing, completed, failed
	SubmittedAt  int64   `json:"submitted_at"`
	CompletedAt  *int64  `json:"completed_at"`
	ErrorMessage *string `json:"error_message"`
}
