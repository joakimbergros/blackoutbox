// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package worker

import (
	"blackoutbox/internal/cups"
	"blackoutbox/internal/monitor"
	"log"
	"time"
)

const (
	checkInterval     = 30 * time.Second
	stuckJobThreshold = 5 * time.Minute
)

type Worker struct {
	monitor *monitor.Monitor
	printer *cups.Printer
	stopCh  chan struct{}
}

func NewWorker(monitor *monitor.Monitor, printer *cups.Printer) *Worker {
	return &Worker{
		monitor: monitor,
		printer: printer,
		stopCh:  make(chan struct{}),
	}
}

func (w *Worker) Start() {
	log.Println("Starting background worker")

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.runChecks()
		case <-w.stopCh:
			log.Println("Stopping background worker")
			return
		}
	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
}

func (w *Worker) runChecks() {
	log.Println("Running periodic checks")

	if err := w.monitor.CheckAllTriggers(); err != nil {
		log.Printf("Error checking triggers: %v", err)
	}

	if err := w.printer.CheckStuckJobs(int(stuckJobThreshold.Seconds())); err != nil {
		log.Printf("Error checking stuck jobs: %v", err)
	}
}
