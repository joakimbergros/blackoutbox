package models

type PrintJob struct {
	Id           int     `json:"id"`
	DocumentId   int     `json:"document_id"`
	CupsJobId    *string `json:"cups_job_id"`
	Status       string  `json:"status"` // pending, printing, completed, failed
	SubmittedAt  int64   `json:"submitted_at"`
	CompletedAt  *int64  `json:"completed_at"`
	ErrorMessage *string `json:"error_message"`
}
