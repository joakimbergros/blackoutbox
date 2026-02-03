package handlers

import (
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
	AddFunc   func(model models.Document) error
	GetFunc   func() ([]models.Document, error)
	LastAdded models.Document
	AddCalled bool
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

func TestDocumentHandlerPost(t *testing.T) {
	tests := []struct {
		name           string
		extId          string
		isSystem       string
		fileContent    []byte
		fileName       string
		expectedStatus int
		expectError    string
	}{
		{
			name:           "Valid request with file and metadata",
			extId:          "doc123",
			isSystem:       "true",
			fileContent:    []byte("test file content"),
			fileName:       "test.txt",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Missing ext_id",
			extId:          "",
			isSystem:       "true",
			fileContent:    []byte("test"),
			fileName:       "test.txt",
			expectedStatus: http.StatusBadRequest,
			expectError:    "ext_id is required",
		},
		{
			name:           "Missing is_system",
			extId:          "doc123",
			isSystem:       "",
			fileContent:    []byte("test"),
			fileName:       "test.txt",
			expectedStatus: http.StatusBadRequest,
			expectError:    "is_system is required",
		},
		{
			name:           "Invalid is_system value",
			extId:          "doc123",
			isSystem:       "not_a_bool",
			fileContent:    []byte("test"),
			fileName:       "test.txt",
			expectedStatus: http.StatusBadRequest,
			expectError:    "is_system needs to be a boolean value",
		},
		{
			name:           "Missing file",
			extId:          "doc123",
			isSystem:       "true",
			fileContent:    nil,
			fileName:       "",
			expectedStatus: http.StatusBadRequest,
			expectError:    "file is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockDocumentStore{}
			handler := &DocumentHandler{Store: mockStore}

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.extId != "" {
				writer.WriteField("ext_id", tt.extId)
			}

			if tt.isSystem != "" {
				writer.WriteField("is_system", tt.isSystem)
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

				if mockStore.LastAdded.ExternalId != "doc123" {
					t.Errorf("Expected ExternalId 'doc123', got '%s'", mockStore.LastAdded.ExternalId)
				}

				if !mockStore.LastAdded.IsSystem {
					t.Error("Expected IsSystem to be true")
				}

				expectedPathPrefix := filepath.Join("uploads", "doc123")
				if mockStore.LastAdded.FilePath[:len(expectedPathPrefix)] != expectedPathPrefix {
					t.Errorf("Expected FilePath to start with '%s', got '%s'", expectedPathPrefix, mockStore.LastAdded.FilePath)
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
