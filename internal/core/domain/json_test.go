package domain

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDependencies struct {
	responseRecorder *httptest.ResponseRecorder
	request          *http.Request
}

func TestWriteJSON(t *testing.T) {
	type args struct {
		status  int
		data    interface{}
		headers []http.Header
	}

	tests := []struct {
		name       string
		setup      func(*testing.T) *testDependencies
		args       args
		want       string
		wantStatus int
		wantHeader http.Header
		wantErr    bool
		cleanup    func(*testing.T, *testDependencies)
	}{
		{
			name: "successful write with no headers",
			setup: func(t *testing.T) *testDependencies {
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			args: args{
				status: http.StatusOK,
				data:   map[string]string{"key": "value"},
			},
			want:       `{"key":"value"}` + "\n",
			wantStatus: http.StatusOK,
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantErr:    false,
			cleanup:    func(t *testing.T, d *testDependencies) {},
		},
		{
			name: "successful write with custom headers",
			setup: func(t *testing.T) *testDependencies {
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			args: args{
				status: http.StatusCreated,
				data:   map[string]string{"key": "value"},
				headers: []http.Header{{
					"X-Custom-Header": {"custom-value"},
				}},
			},
			want:       `{"key":"value"}` + "\n",
			wantStatus: http.StatusCreated,
			wantHeader: http.Header{
				"Content-Type":    []string{"application/json"},
				"X-Custom-Header": []string{"custom-value"},
			},
			wantErr: false,
			cleanup: func(t *testing.T, d *testDependencies) {},
		},
		{
			name: "error on invalid JSON",
			setup: func(t *testing.T) *testDependencies {
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			args: args{
				status: http.StatusOK,
				data:   make(chan int), // Channels can't be JSON marshaled
			},
			wantErr: true,
			cleanup: func(t *testing.T, d *testDependencies) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			deps := tt.setup(t)
			defer tt.cleanup(t, deps)

			// Execute
			err := WriteJSON(deps.responseRecorder, tt.args.status, tt.args.data, tt.args.headers...)

			// Verify
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.wantStatus, deps.responseRecorder.Code)

			// Check headers
			for k, v := range tt.wantHeader {
				assert.Equal(t, v, deps.responseRecorder.Header()[k])
			}

			// Check body
			var actual, expected interface{}
			err = json.Unmarshal(deps.responseRecorder.Body.Bytes(), &actual)
			require.NoError(t, err)
			err = json.Unmarshal([]byte(tt.want), &expected)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestWriteJSONError(t *testing.T) {
	type args struct {
		status  int
		message string
		headers []http.Header
	}

	tests := []struct {
		name        string
		setup       func(*testing.T) *testDependencies
		args        args
		want        *JSONResponse
		wantStatus  int
		wantErr     bool
		wantHeaders http.Header
		cleanup     func(*testing.T, *testDependencies)
	}{
		{
			name: "successful error response",
			setup: func(t *testing.T) *testDependencies {
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			args: args{
				status:  http.StatusBadRequest,
				message: "invalid input",
				headers: []http.Header{{"X-Request-ID": {"123"}}},
			},
			want: &JSONResponse{
				Error:   true,
				Message: "invalid input",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
			wantHeaders: http.Header{
				"Content-Type": {"application/json"},
				"X-Request-ID": {"123"},
			},
			cleanup: func(t *testing.T, d *testDependencies) {},
		},
		{
			name: "empty error message",
			setup: func(t *testing.T) *testDependencies {
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			args: args{
				status:  http.StatusInternalServerError,
				message: "",
			},
			want: &JSONResponse{
				Error:   true,
				Message: "",
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
			wantHeaders: http.Header{
				"Content-Type": {"application/json"},
			},
			cleanup: func(t *testing.T, d *testDependencies) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			deps := tt.setup(t)
			defer tt.cleanup(t, deps)

			// Execute
			err := WriteJSONError(deps.responseRecorder, tt.args.status, tt.args.message, tt.args.headers...)

			// Verify
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.wantStatus, deps.responseRecorder.Code)

			// Check response body matches expected
			var resp JSONResponse
			err = json.Unmarshal(deps.responseRecorder.Body.Bytes(), &resp)
			require.NoError(t, err)
			assert.Equal(t, tt.want, &resp)
		})
	}
}

func TestReadJSON(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	type args struct {
		target any
		req    *http.Request
	}

	tests := []struct {
		name        string
		setup       func(*testing.T) *testDependencies
		args        args
		want        any
		wantErr     bool
		errContains string
		checkHeader bool
		cleanup     func(*testing.T, *testDependencies)
	}{
		{
			name: "successful read",
			setup: func(t *testing.T) *testDependencies {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"test","value":42}`))
				req.Header.Set("Content-Type", "application/json")
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			args: args{
				target: &testStruct{},
			},
			want: &testStruct{
				Name:  "test",
				Value: 42,
			},
			wantErr:     false,
			checkHeader: false,
			cleanup:     func(t *testing.T, d *testDependencies) {},
		},
		{
			name: "empty body",
			setup: func(t *testing.T) *testDependencies {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))
				req.Header.Set("Content-Type", "application/json")
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			args: args{
				target: &testStruct{},
			},
			want:        nil,
			wantErr:     true,
			errContains: "EOF",
			checkHeader: false,
			cleanup:     func(t *testing.T, d *testDependencies) {},
		},
		{
			name: "invalid JSON",
			setup: func(t *testing.T) *testDependencies {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("{invalid-json}"))
				req.Header.Set("Content-Type", "application/json")
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			args: args{
				target: &testStruct{},
			},
			want:        nil,
			wantErr:     true,
			errContains: "invalid character 'i' looking for beginning of object key string",
			checkHeader: false,
			cleanup:     func(t *testing.T, d *testDependencies) {},
		},
		{
			name: "unknown field",
			setup: func(t *testing.T) *testDependencies {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"test","value":42,"unknown":"field"}`))
				req.Header.Set("Content-Type", "application/json")
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			args: args{
				target: &testStruct{},
			},
			want:        nil,
			wantErr:     true,
			errContains: "json: unknown field",
			checkHeader: false,
			cleanup:     func(t *testing.T, d *testDependencies) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			deps := tt.setup(t)
			defer tt.cleanup(t, deps)

			// Execute
			err := ReadJSON(deps.responseRecorder, deps.request, tt.args.target)

			// Verify
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, tt.args.target)

			if tt.checkHeader {
				assert.Equal(t, "application/json", deps.responseRecorder.Header().Get("Content-Type"))
			}
		})
	}
}

func TestReadJSON_MaxBytes(t *testing.T) {
	type args struct {
		target any
		req    *http.Request
	}

	tests := []struct {
		name        string
		setup       func(*testing.T) *testDependencies
		args        args
		wantErr     bool
		errContains string
		cleanup     func(*testing.T, *testDependencies)
	}{
		{
			name: "request body too large",
			setup: func(t *testing.T) *testDependencies {
				// Create a request with body larger than maxBytes (1MB)
				hugeBody := make([]byte, 2*1024*1024) // 2MB
				req := httptest.NewRequest("POST", "/test", bytes.NewReader(hugeBody))
				req.Header.Set("Content-Type", "application/json")
				return &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			args: args{
				target: &map[string]any{},
			},
			wantErr:     true,
			errContains: "invalid character '\\x00' looking for beginning of value",
			cleanup:     func(t *testing.T, d *testDependencies) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			deps := tt.setup(t)
			defer tt.cleanup(t, deps)

			// Execute
			err := ReadJSON(deps.responseRecorder, deps.request, tt.args.target)

			// Verify
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
		})
	}
}
