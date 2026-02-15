// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package monitor

import (
	"blackoutbox/internal/models"
	"blackoutbox/internal/stores"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	maxRetries     = 3
	checkTimeout   = 10 * time.Second
	triggeredState = "triggered"
)

type Monitor struct {
	triggerStore    stores.TriggerStoreInterface
	documentStore   stores.DocumentStoreInterface
	templateStore   stores.TemplateStoreInterface
	printJobStore   stores.PrintJobStoreInterface
	printJobCreator PrintJobCreator
}

type PrintJobCreator interface {
	CreatePrintJob(documentId int64, filePath string) error
}

func NewMonitor(
	triggerStore stores.TriggerStoreInterface,
	documentStore stores.DocumentStoreInterface,
	templateStore stores.TemplateStoreInterface,
	printJobStore stores.PrintJobStoreInterface,
	printJobCreator PrintJobCreator,
) *Monitor {
	return &Monitor{
		triggerStore:    triggerStore,
		documentStore:   documentStore,
		templateStore:   templateStore,
		printJobStore:   printJobStore,
		printJobCreator: printJobCreator,
	}
}

func (m *Monitor) CheckTrigger(trigger models.Trigger) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("HEAD", trigger.Url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	startTime := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	now := time.Now().Unix()

	if err != nil {
		return m.handleFailure(trigger, now, fmt.Sprintf("connection error: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return m.handleSuccess(trigger, now)
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return m.handleFailure(trigger, now, fmt.Sprintf("client error: %s", resp.Status))
	}

	if resp.StatusCode >= 500 {
		return m.handleFailure(trigger, now, fmt.Sprintf("server error: %s", resp.Status))
	}

	if duration > 5*time.Second {
		return m.handleFailure(trigger, now, fmt.Sprintf("slow response: %v", duration))
	}

	return m.handleFailure(trigger, now, fmt.Sprintf("unexpected status: %s", resp.Status))
}

func (m *Monitor) handleFailure(trigger models.Trigger, now int64, reason string) error {
	log.Printf("Trigger %d failed: %s", trigger.Id, reason)

	triggerId := trigger.Id

	if trigger.LastFailedAt == nil {
		nowCopy := now
		trigger.LastFailedAt = &nowCopy
		trigger.Status = "error"
		trigger.RetryCount = 1
		return m.triggerStore.Update(trigger)
	}

	trigger.RetryCount++

	if trigger.RetryCount >= maxRetries {
		failureDuration := now - *trigger.LastFailedAt
		if failureDuration >= int64(trigger.BufferSeconds) {
			log.Printf("Trigger %d buffer exceeded (%ds), triggering print jobs", trigger.Id, failureDuration)
			trigger.Status = triggeredState
			if err := m.triggerStore.UpdateStatus(triggerId, triggeredState); err != nil {
				return fmt.Errorf("failed to update trigger status: %w", err)
			}
			return m.triggerPrintJobs(trigger.SystemId)
		}
	}

	return m.triggerStore.Update(trigger)
}

func (m *Monitor) handleSuccess(trigger models.Trigger, now int64) error {
	log.Printf("Trigger %d check successful", trigger.Id)

	triggerId := trigger.Id

	if trigger.RetryCount > 0 || trigger.Status != "ok" {
		return m.triggerStore.ResetRetryCount(triggerId)
	}

	nowCopy := now
	trigger.LastCheckedAt = &nowCopy
	return m.triggerStore.Update(trigger)
}

func (m *Monitor) triggerPrintJobs(systemId int64) error {
	documents, err := m.documentStore.GetBySystemId(systemId)
	if err != nil {
		return fmt.Errorf("failed to get documents for system %s: %w", systemId, err)
	}

	for _, doc := range documents {
		templates, err := m.templateStore.GetByFileReference(doc.FileReference)
		if err != nil {
			log.Printf("Failed to gather templates from db for document %d: %v", doc.Id, err)
		}

		if err := m.printJobCreator.CreatePrintJob(doc.Id, doc.FilePath); err != nil {
			log.Printf("Failed to create print job for document %d: %v", doc.Id, err)
		}

		//TODO Should support multiple templates tied to single file_id?
		if templates != nil {
			if err := m.printJobCreator.CreatePrintJob(templates.Id, templates.FilePath); err != nil {
				log.Printf("Failed to create print job for document %d: %v", doc.Id, err)
			}
		}
	}

	return nil
}

func (m *Monitor) CheckAllTriggers() error {
	triggers, err := m.triggerStore.Get()
	if err != nil {
		return fmt.Errorf("failed to get triggers: %w", err)
	}

	for _, trigger := range triggers {
		if trigger.Status == triggeredState {
			continue
		}

		now := time.Now().Unix()
		nowCopy := now
		trigger.LastCheckedAt = &nowCopy

		if err := m.CheckTrigger(trigger); err != nil {
			log.Printf("Error checking trigger %d: %v", trigger.Id, err)
		}
	}

	return nil
}
