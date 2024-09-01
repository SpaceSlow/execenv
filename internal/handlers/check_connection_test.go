package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SpaceSlow/execenv/internal/storages"
)

var _ storages.ICheckConnection = (*MockCheckStorage)(nil)

type MockCheckStorage struct {
	returnCheck bool
}

func (s *MockCheckStorage) CheckConnection() bool {
	return s.returnCheck
}

func TestDBHandler_Ping(t *testing.T) {
	tests := []struct {
		storage        storages.ICheckConnection
		name           string
		expectedStatus int
	}{
		{
			name:           "CheckConnection() return true",
			storage:        &MockCheckStorage{returnCheck: true},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "CheckConnection() return false",
			storage:        &MockCheckStorage{returnCheck: false},
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CheckConnectionHandler{storage: tt.storage}
			req, err := http.NewRequest("GET", "/ping", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.Ping(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
