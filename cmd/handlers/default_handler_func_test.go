package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	type fields struct {
		method string
		url    string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "GET method incorrect metric type",
			fields: fields{
				method: http.MethodGet,
				url:    "/update/unknown_type/metric/42",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "POST method incorrect metric type",
			fields: fields{
				method: http.MethodPost,
				url:    "/update/unknown_type/metric/42",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "empty metric type",
			fields: fields{
				method: http.MethodPost,
				url:    "/update/",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "root path",
			fields: fields{
				method: http.MethodPost,
				url:    "/",
			},
			want: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.fields.method, test.fields.url, nil)

			w := httptest.NewRecorder()
			DefaultHandlerFunc(w, request)
			res := w.Result()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}
