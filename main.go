package main

import (
	"blackoutbox/internal/cups"
	"blackoutbox/internal/handlers/documents"
	"blackoutbox/internal/handlers/printjobs"
	"blackoutbox/internal/handlers/systems"
	"blackoutbox/internal/handlers/triggers"
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

	triggerStore := stores.TriggerStore{Db: db}
	triggerHandler := triggers.TriggerHandler{Store: &triggerStore}

	printJobStore := stores.PrintJobStore{Db: db}
	printJobHandler := printjobs.PrintJobHandler{Store: &printJobStore}

	systemStore := stores.SystemStore{
		Db:        db,
		FilesRoot: "uploads",
	}

	systemHandler := systems.SystemHandler{
		SystemStore: &systemStore,
		UploadRoot:  "uploads",
	}

	printer := cups.NewPrinter(&printJobStore)
	monitorService := monitor.NewMonitor(&triggerStore, &documentStore, &printJobStore, printer)
	workerService := worker.NewWorker(monitorService, printer)

	go workerService.Start()

	mux := http.NewServeMux()

	mux.Handle("GET /documents", documentHandler.Get())
	mux.Handle("GET /documents/{id}", documentHandler.GetById())
	mux.Handle("POST /documents", documentHandler.Post())
	mux.Handle("PATCH /documents", documentHandler.Update())

	mux.Handle("GET /triggers", triggerHandler.Get())
	mux.Handle("GET /triggers/{id}", triggerHandler.GetById())
	mux.Handle("POST /triggers", triggerHandler.Post())
	mux.Handle("DELETE /triggers/{id}", triggerHandler.Delete())

	mux.Handle("GET /print_jobs", printJobHandler.Get())
	mux.Handle("GET /print_jobs/{id}", printJobHandler.GetById())
	mux.Handle("GET /print_jobs/stuck", printJobHandler.GetStuck())

	// System routes (bulk / orchestration)
	mux.Handle("POST /systems/{system_id}/sync", systemHandler.Sync())
	mux.Handle("DELETE /systems/{system_id}", systemHandler.Delete())

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
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
