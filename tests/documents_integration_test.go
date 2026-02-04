// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tests

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"blackoutbox/internal/handlers/documents"
	"blackoutbox/internal/storage"
	"blackoutbox/internal/stores"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create the documents table
	_, err = db.Exec(`
		CREATE TABLE documents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			system_id TEXT NOT NULL,
			file_id TEXT NOT NULL,
			file_path TEXT NOT NULL,
			print_at INTEGER,
			last_printed_at INTEGER,
			tags TEXT,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create documents table: %v", err)
	}

	return db
}

func createMultipartRequest(t *testing.T, systemId, fileId, tags, printAt string, fileContent []byte, fileName string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("system_id", systemId)
	writer.WriteField("file_id", fileId)

	if tags != "" {
		writer.WriteField("tags", tags)
	}

	if printAt != "" {
		writer.WriteField("print_at", printAt)
	}

	if fileContent != nil {
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		part.Write(fileContent)
	}

	writer.Close()

	req := httptest.NewRequest("POST", "/documents", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func cleanupUploads(filePath string) {
	if filePath != "" {
		os.Remove(filePath)
		dir := filepath.Dir(filePath)
		if dir != storage.DocumentsRoot && dir != "." {
			os.Remove(dir)
			parentDir := filepath.Dir(dir)
			if parentDir == storage.DocumentsRoot {
				os.Remove(parentDir)
			}
		}
	}
}

func TestIntegrationDocumentHandlerPost(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	documentStore := stores.DocumentStore{Db: db}
	documentHandler := documents.DocumentHandler{Store: &documentStore}

	tests := []struct {
		name           string
		request        *http.Request
		expectedStatus int
		expectError    string
		cleanupPath    string
	}{
		{
			name:           "Valid document upload with all fields",
			request:        createMultipartRequest(t, "system123", "file456", `["daily", "instant"]`, "1738581234", []byte("test file content"), "test.txt"),
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Valid document upload with minimal fields",
			request:        createMultipartRequest(t, "system456", "file789", "", "", []byte("minimal content"), "minimal.txt"),
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Missing system_id",
			request:        createMultipartRequest(t, "", "file456", "", "", []byte("test"), "test.txt"),
			expectedStatus: http.StatusBadRequest,
			expectError:    "system_id is required",
		},
		{
			name:           "Missing file_id",
			request:        createMultipartRequest(t, "system123", "", "", "", []byte("test"), "test.txt"),
			expectedStatus: http.StatusBadRequest,
			expectError:    "file_id is required",
		},
		{
			name:           "Invalid print_at timestamp",
			request:        createMultipartRequest(t, "system123", "file456", "", "not_a_timestamp", []byte("test"), "test.txt"),
			expectedStatus: http.StatusBadRequest,
			expectError:    "print_at must be a valid Unix timestamp",
		},
		{
			name:           "Invalid tags JSON",
			request:        createMultipartRequest(t, "system123", "file456", "not_json_array", "", []byte("test"), "test.txt"),
			expectedStatus: http.StatusBadRequest,
			expectError:    "tags must be a valid JSON array",
		},
		{
			name:           "Missing file",
			request:        createMultipartRequest(t, "system123", "file456", "", "", nil, ""),
			expectedStatus: http.StatusBadRequest,
			expectError:    "file is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			postHandler := documentHandler.Post()
			postHandler(rr, tt.request)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectError != "" {
				bodyStr := rr.Body.String()
				if len(bodyStr) > 0 && bodyStr[len(bodyStr)-1] == '\n' {
					bodyStr = bodyStr[:len(bodyStr)-1]
				}
				if bodyStr != tt.expectError {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectError, bodyStr)
				}
			}

			// Verify document was stored in database for successful requests
			if tt.expectedStatus == http.StatusCreated {
				rows, err := db.Query("SELECT system_id, file_id, tags FROM documents WHERE system_id = ?", "system123")
				if err != nil {
					t.Fatalf("Failed to query documents: %v", err)
				}
				defer rows.Close()

				if !rows.Next() {
					t.Error("Expected document to be stored in database")
				}

				var systemId, fileId, tagsJSON string
				err = rows.Scan(&systemId, &fileId, &tagsJSON)
				if err != nil {
					t.Fatalf("Failed to scan document: %v", err)
				}

				if systemId != "system123" {
					t.Errorf("Expected system_id 'system123', got '%s'", systemId)
				}
				if fileId != "file456" {
					t.Errorf("Expected file_id 'file456', got '%s'", fileId)
				}

				// Clean up uploaded file
				cleanupUploads(tt.cleanupPath)
			}
		})
	}
}

func TestIntegrationDocumentHandlerGet(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	documentStore := stores.DocumentStore{Db: db}
	documentHandler := documents.DocumentHandler{Store: &documentStore}

	// First, insert a document via POST
	postReq := createMultipartRequest(t, "system123", "file456", `["test"]`, "", []byte("test content"), "test.txt")
	postRR := httptest.NewRecorder()
	postHandler := documentHandler.Post()
	postHandler(postRR, postReq)

	if postRR.Code != http.StatusCreated {
		t.Fatalf("Failed to create test document: status %d", postRR.Code)
	}

	// Now test GET
	getReq := httptest.NewRequest("GET", "/documents", nil)
	getRR := httptest.NewRecorder()

	getHandler := documentHandler.Get()
	getHandler(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", getRR.Code)
	}

	contentType := getRR.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	body := getRR.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}

	// Clean up
	cleanupUploads("")
}

func TestIntegrationDocumentHandlerLargeFile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	documentStore := stores.DocumentStore{Db: db}
	documentHandler := documents.DocumentHandler{Store: &documentStore}

	// Create a file larger than 10MB limit
	largeContent := make([]byte, 11<<20) // 11MB
	req := createMultipartRequest(t, "system123", "file456", "", "", largeContent, "large.txt")

	rr := httptest.NewRecorder()
	postHandler := documentHandler.Post()
	postHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for large file, got %d", rr.Code)
	}

	bodyStr := rr.Body.String()
	if len(bodyStr) > 0 && bodyStr[len(bodyStr)-1] == '\n' {
		bodyStr = bodyStr[:len(bodyStr)-1]
	}
	expectedError := "File size exceeds 10MB limit"
	if bodyStr != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, bodyStr)
	}
}

func TestIntegrationDocumentHandlerMultipleUploads(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	documentStore := stores.DocumentStore{Db: db}
	documentHandler := documents.DocumentHandler{Store: &documentStore}

	numUploads := 5

	// Upload multiple documents sequentially
	for i := range numUploads {
		systemId := fmt.Sprintf("system%d", i)
		fileId := fmt.Sprintf("file%d", i)
		req := createMultipartRequest(t, systemId, fileId, "", "", []byte("test content"), "test.txt")
		rr := httptest.NewRecorder()
		postHandler := documentHandler.Post()
		postHandler(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Upload %d failed with status %d", i, rr.Code)
		}
	}

	// Verify all documents are in database
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM documents").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count documents: %v", err)
	}

	if count != numUploads {
		t.Errorf("Expected %d documents in database, got %d", numUploads, count)
	}

	// Verify GET returns all documents
	getReq := httptest.NewRequest("GET", "/documents", nil)
	getRR := httptest.NewRecorder()
	getHandler := documentHandler.Get()
	getHandler(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Errorf("Expected status 200 for GET, got %d", getRR.Code)
	}

	// Clean up
	cleanupUploads("")
}
