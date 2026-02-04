// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cups

import (
	"blackoutbox/internal/models"
	"blackoutbox/internal/stores"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Printer struct {
	printJobStore stores.PrintJobStoreInterface
}

func NewPrinter(printJobStore stores.PrintJobStoreInterface) *Printer {
	return &Printer{
		printJobStore: printJobStore,
	}
}

func (p *Printer) CreatePrintJob(documentId int, filePath string) error {
	jobId, err := p.submitPrint(filePath)
	if err != nil {
		p.recordFailedJob(documentId, err.Error())
		return fmt.Errorf("failed to submit print job: %w", err)
	}

	return p.recordSuccessfulJob(documentId, jobId)
}

func (p *Printer) submitPrint(filePath string) (string, error) {
	cmd := exec.Command("lp", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("lp command failed: %w, output: %s", err, string(output))
	}

	jobId := p.parseJobId(string(output))
	if jobId == "" {
		return "", fmt.Errorf("could not parse job id from lp output: %s", string(output))
	}

	log.Printf("Submitted print job %s for file: %s", jobId, filePath)
	return jobId, nil
}

func (p *Printer) parseJobId(output string) string {
	re := regexp.MustCompile(`request id is \S+-(\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func (p *Printer) recordSuccessfulJob(documentId int, cupsJobId string) error {
	now := time.Now().Unix()

	job := models.PrintJob{
		DocumentId:  documentId,
		CupsJobId:   &cupsJobId,
		Status:      "printing",
		SubmittedAt: now,
	}

	if err := p.printJobStore.Add(job); err != nil {
		return fmt.Errorf("failed to record print job: %w", err)
	}

	return nil
}

func (p *Printer) recordFailedJob(documentId int, errorMessage string) error {
	now := time.Now().Unix()

	job := models.PrintJob{
		DocumentId:   documentId,
		Status:       "failed",
		SubmittedAt:  now,
		ErrorMessage: &errorMessage,
	}

	if err := p.printJobStore.Add(job); err != nil {
		return fmt.Errorf("failed to record failed print job: %w", err)
	}

	return nil
}

func (p *Printer) CheckJobStatus(cupsJobId string) (string, error) {
	cmd := exec.Command("lpq")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("lpq command failed: %w, output: %s", err, string(output))
	}

	return p.parseJobStatus(string(output), cupsJobId)
}

func (p *Printer) parseJobStatus(output, cupsJobId string) (string, error) {
	lines := strings.SplitSeq(output, "\n")
	for line := range lines {
		if strings.Contains(line, cupsJobId) {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				status := strings.ToLower(fields[3])
				if status == "active" || status == "printing" {
					return "printing", nil
				}
				if status == "completed" || status == "done" {
					return "completed", nil
				}
				if status == "held" || status == "error" {
					return "failed", nil
				}
			}
		}
	}

	return "completed", nil
}

func (p *Printer) UpdateJobStatus(jobId string) error {
	job, err := p.printJobStore.GetById(jobId)
	if err != nil {
		return fmt.Errorf("failed to get print job: %w", err)
	}

	if job.CupsJobId == nil {
		return fmt.Errorf("job has no cups job id")
	}

	status, err := p.CheckJobStatus(*job.CupsJobId)
	if err != nil {
		return fmt.Errorf("failed to check job status: %w", err)
	}

	if status != job.Status {
		if status == "completed" {
			now := time.Now().Unix()
			job.Status = status
			job.CompletedAt = &now
		} else {
			job.Status = status
		}

		return p.printJobStore.Update(*job)
	}

	return nil
}

func (p *Printer) CheckStuckJobs(thresholdSeconds int) error {
	stuckJobs, err := p.printJobStore.GetStuckJobs(thresholdSeconds)
	if err != nil {
		return fmt.Errorf("failed to get stuck jobs: %w", err)
	}

	for _, job := range stuckJobs {
		log.Printf("Found stuck job %d (document %d), status: %s", job.Id, job.DocumentId, job.Status)
		if job.CupsJobId != nil {
			if err := p.UpdateJobStatus(strconv.Itoa(job.Id)); err != nil {
				log.Printf("Failed to update job %d status: %v", job.Id, err)
			}
		}
	}

	return nil
}
