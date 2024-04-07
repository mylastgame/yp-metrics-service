package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:   "test /",
			method: http.MethodGet,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
			request: "/",
		},
		{
			name:   "wrong URI /updater",
			method: http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
			request: "/updater",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(DefaultHandler())
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
