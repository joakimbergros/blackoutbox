package tests

import (
	"blackoutbox/internal/handlers/documents"
	"blackoutbox/internal/models"
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type MockDocumentStore struct {
	AddFunc     func(model models.Document) error
	GetFunc     func() ([]models.Document, error)
	UpdateFunc  func(model models.Document) error
	GetByIdFunc func(id int) (*models.Document, error)
	LastAdded   models.Document
	AddCalled   bool
}

func (m *MockDocumentStore) Add(model models.Document) error {
	m.LastAdded = model
	m.AddCalled = true
	if m.AddFunc != nil {
		return m.AddFunc(model)
	}
	return nil
}

func (m *MockDocumentStore) Get() ([]models.Document, error) {
	if m.GetFunc != nil {
		return m.GetFunc()
	}
	return []models.Document{}, nil
}

func (m *MockDocumentStore) Update(model models.Document) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(model)
	}
	return nil
}

func (m *MockDocumentStore) GetById(id int) (*models.Document, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(id)
	}
	return nil, nil
}

func TestDocumentHandlerPost(t *testing.T) {
	tests := []struct {
		name           string
		systemId       string
		fileId         string
		tags           string
		printAt        string
		fileContent    []byte
		fileName       string
		expectedStatus int
		expectError    string
	}{
		{
			name:           "Valid request with all fields",
			systemId:       "system123",
			fileId:         "file456",
			tags:           `["daily", "instant"]`,
			printAt:        "1738581234",
			fileContent:    []byte("test file content"),
			fileName:       "test.txt",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Valid request with minimal fields",
			systemId:       "system123",
			fileId:         "file456",
			fileContent:    []byte("test"),
			fileName:       "test.txt",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Missing system_id",
			systemId:       "",
			fileId:         "file456",
			fileContent:    []byte("test"),
			fileName:       "test.txt",
			expectedStatus: http.StatusBadRequest,
			expectError:    "system_id is required",
		},
		{
			name:           "Missing file_id",
			systemId:       "system123",
			fileId:         "",
			fileContent:    []byte("test"),
			fileName:       "test.txt",
			expectedStatus: http.StatusBadRequest,
			expectError:    "file_id is required",
		},
		{
			name:           "Invalid print_at (not a number)",
			systemId:       "system123",
			fileId:         "file456",
			printAt:        "not_a_timestamp",
			fileContent:    []byte("test"),
			fileName:       "test.txt",
			expectedStatus: http.StatusBadRequest,
			expectError:    "print_at must be a valid Unix timestamp",
		},
		{
			name:           "Invalid tags JSON",
			systemId:       "system123",
			fileId:         "file456",
			tags:           "not_json_array",
			fileContent:    []byte("test"),
			fileName:       "test.txt",
			expectedStatus: http.StatusBadRequest,
			expectError:    "tags must be a valid JSON array",
		},
		{
			name:           "Missing file",
			systemId:       "system123",
			fileId:         "file456",
			fileContent:    nil,
			fileName:       "",
			expectedStatus: http.StatusBadRequest,
			expectError:    "file is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDocumentStore{}
			handler := &documents.DocumentHandler{Store: mockStore}

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.systemId != "" {
				writer.WriteField("system_id", tt.systemId)
			}

			if tt.fileId != "" {
				writer.WriteField("file_id", tt.fileId)
			}

			if tt.tags != "" {
				writer.WriteField("tags", tt.tags)
			}

			if tt.printAt != "" {
				writer.WriteField("print_at", tt.printAt)
			}

			if tt.fileContent != nil {
				part, err := writer.CreateFormFile("file", tt.fileName)
				if err != nil {
					t.Fatalf("Failed to create form file: %v", err)
				}
				part.Write(tt.fileContent)
			}

			writer.Close()

			req := httptest.NewRequest("POST", "/documents", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rr := httptest.NewRecorder()

			postHandler := handler.Post()
			postHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectError != "" {
				bodyStr := rr.Body.String()
				if bodyStr[:len(bodyStr)-1] != tt.expectError {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectError, bodyStr[:len(bodyStr)-1])
				}
			}

			if tt.expectedStatus == http.StatusCreated {
				if !mockStore.AddCalled {
					t.Error("Expected Store.Add to be called")
				}

				if mockStore.LastAdded.SystemId != "system123" {
					t.Errorf("Expected SystemId 'system123', got '%s'", mockStore.LastAdded.SystemId)
				}

				if mockStore.LastAdded.FileId != "file456" {
					t.Errorf("Expected FileId 'file456', got '%s'", mockStore.LastAdded.FileId)
				}

				expectedPathPrefix := filepath.Join("uploads", "system123")
				if mockStore.LastAdded.FilePath[:len(expectedPathPrefix)] != expectedPathPrefix {
					t.Errorf("Expected FilePath to start with '%s', got '%s'", expectedPathPrefix, mockStore.LastAdded.FilePath)
				}

				if tt.tags != "" {
					if len(mockStore.LastAdded.Tags) != 2 {
						t.Errorf("Expected 2 tags, got %d", len(mockStore.LastAdded.Tags))
					}
					if mockStore.LastAdded.Tags[0] != "daily" || mockStore.LastAdded.Tags[1] != "instant" {
						t.Errorf("Expected tags [daily, instant], got %v", mockStore.LastAdded.Tags)
					}
				}

				if tt.printAt != "" {
					if mockStore.LastAdded.PrintAt == nil {
						t.Error("Expected PrintAt to be set")
					} else if *mockStore.LastAdded.PrintAt != 1738581234 {
						t.Errorf("Expected PrintAt 1738581234, got %d", *mockStore.LastAdded.PrintAt)
					}
				}

				// Clean up uploaded file
				if _, err := os.Stat(mockStore.LastAdded.FilePath); err == nil {
					os.Remove(mockStore.LastAdded.FilePath)
					os.Remove(filepath.Dir(mockStore.LastAdded.FilePath))
				}
			}
		})
	}
}
