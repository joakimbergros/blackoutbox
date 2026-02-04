// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"blackoutbox/internal/cups"
	"blackoutbox/internal/handlers/documents"
	"blackoutbox/internal/handlers/printjobs"
	"blackoutbox/internal/handlers/systems"
	"blackoutbox/internal/handlers/templates"
	"blackoutbox/internal/handlers/triggers"
	"blackoutbox/internal/middleware"
	"blackoutbox/internal/monitor"
	"blackoutbox/internal/stores"
	"blackoutbox/internal/worker"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:./app.db?_foreign_keys=on")
	if err != nil {
		log.Panic("unable to connect to db")
	}

	documentStore := stores.DocumentStore{Db: db}
	documentHandler := documents.DocumentHandler{Store: &documentStore}

	templateStore := stores.TemplateStore{Db: db}
	templateHandler := templates.TemplatesHandler{Store: &templateStore}

	triggerStore := stores.TriggerStore{Db: db}
	triggerHandler := triggers.TriggerHandler{Store: &triggerStore}

	printJobStore := stores.PrintJobStore{Db: db}
	printJobHandler := printjobs.PrintJobHandler{Store: &printJobStore}

	systemStore := stores.SystemStore{
		Db: db,
	}

	systemHandler := systems.SystemHandler{
		SystemStore: &systemStore,
	}

	printer := cups.NewPrinter(&printJobStore)
	monitorService := monitor.NewMonitor(&triggerStore, &documentStore, &templateStore, &printJobStore, printer)
	workerService := worker.NewWorker(monitorService, printer)

	go workerService.Start()

	baseMiddleware := middleware.Chain{middleware.LogRequest}
	authMiddleware := append(baseMiddleware, middleware.AuthMiddleware)

	mux := http.NewServeMux()

	mux.Handle("GET /documents", baseMiddleware.Then(documentHandler.Get()))
	mux.Handle("GET /documents/{id}", baseMiddleware.Then(documentHandler.GetById()))
	mux.Handle("POST /documents", authMiddleware.Then(documentHandler.Post()))
	mux.Handle("PATCH /documents", authMiddleware.Then(documentHandler.Update()))

	mux.Handle("GET /templates", authMiddleware.Then(templateHandler.Get()))
	mux.Handle("POST /templates", authMiddleware.Then(templateHandler.Post()))
	mux.Handle("DELETE /templates", authMiddleware.Then(templateHandler.Delete()))

	mux.Handle("GET /triggers", baseMiddleware.Then(triggerHandler.Get()))
	mux.Handle("GET /triggers/{id}", baseMiddleware.Then(triggerHandler.GetById()))
	mux.Handle("POST /triggers", authMiddleware.Then(triggerHandler.Post()))
	mux.Handle("DELETE /triggers/{id}", authMiddleware.Then(triggerHandler.Delete()))

	mux.Handle("GET /print_jobs", baseMiddleware.Then(printJobHandler.Get()))
	mux.Handle("GET /print_jobs/{id}", baseMiddleware.Then(printJobHandler.GetById()))
	mux.Handle("GET /print_jobs/stuck", baseMiddleware.Then(printJobHandler.GetStuck()))

	// System routes (bulk / orchestration)
	mux.Handle("POST /systems/{system_id}/sync", authMiddleware.Then(systemHandler.Sync()))
	mux.Handle("DELETE /systems/{system_id}", authMiddleware.Then(systemHandler.Delete()))

	serverTimeout := 5 * time.Second
	server := &http.Server{
		Addr:              ":3000",
		Handler:           mux,
		ReadHeaderTimeout: serverTimeout,
		ReadTimeout:       serverTimeout,
	}

	go func() {
		log.Println("Server starting on :3000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("Unable to start server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	workerService.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
