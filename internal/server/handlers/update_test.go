package handlers

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
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
				statusCode:  http.StatusMethodNotAllowed,
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
		{
			name:   "wrong type",
			method: http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
			request: "/update/ggauge/g1/2",
		},
		{
			name:   "empty gauge title /update/gauge/",
			method: http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusNotFound,
			},
			request: "/update/gauge/",
		},
		{
			name:   "empty counter title /update/counter/",
			method: http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusNotFound,
			},
			request: "/update/counter/",
		},
		{
			name:   "add gauge g1",
			method: http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusOK,
			},
			request: "/update/gauge/g1/0.22",
		},
		{
			name:   "add wrong gauge g2",
			method: http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
			request: "/update/gauge/g2/0.2d2",
		},
		{
			name:   "add counter c1",
			method: http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusOK,
			},
			request: "/update/counter/c1/2",
		},
		{
			name:   "add wrong counter c1",
			method: http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
			request: "/update/counter/c1/0.22",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()

			App := app.New()
			h := http.HandlerFunc(UpdateHandler(App))
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
