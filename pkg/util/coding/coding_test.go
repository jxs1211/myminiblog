package coding

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleHandler(t *testing.T) {
	tt := []struct {
		name         string
		value        string
		expectedCode int
		res          int
		message      string
	}{
		{name: "empty", value: "", expectedCode: 400, message: "no value"},
		{name: "invalid", value: "abc", expectedCode: 400, message: "invalid value"},
		{name: "1", value: "1", expectedCode: 200, res: 2},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://localhost:8080/double?v="+tc.value, nil)
			assert.Nil(t, err)
			rec := httptest.NewRecorder()
			doubleHandler(rec, req)

			res := rec.Result()
			assert.NotNil(t, res)
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedCode, res.StatusCode)
			msg := strings.TrimSpace(string(data))

			if tc.message != "" {
				assert.Equal(t, tc.message, msg)
				return
			}
			i, err := strconv.Atoi(msg)
			assert.Nil(t, err)
			assert.Equal(t, tc.res, i)
		})
	}
}

func TestHandler(t *testing.T) {

}

func Test_handler(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		expectErr          string
		expectedStatusCode int
		wantRes            int
	}{
		{name: "Test_handler", path: "/double?v=1", expectedStatusCode: 200, wantRes: 2},
		{name: "not value", path: "/double", expectedStatusCode: 400, expectErr: "no value"},
		{name: "invalid value", path: "/double?v=abc", expectedStatusCode: 400, expectErr: "invalid value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handler()
			assert.NotNil(t, got)
			server := httptest.NewServer(got)
			url := server.URL + tt.path
			res, err := http.Get(url)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			assert.Nil(t, err)

			if tt.expectErr != "" {
				assert.Equal(t, tt.expectErr, strings.TrimSpace(string(data)))
				return
			}
			gotRes, err := strconv.Atoi(strings.TrimSpace(string(data)))
			assert.Nil(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}
