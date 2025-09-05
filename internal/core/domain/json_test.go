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

type writeJSONTestCase struct {
	name   string
	deps   *testDependencies
	before func(t *testing.T, tc *writeJSONTestCase)
	after  func(t *testing.T, tc *writeJSONTestCase)
	args   struct {
		status  int
		data    any
		headers []http.Header
	}
	expected struct {
		status  int
		headers http.Header
		body    string
		err     bool
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []writeJSONTestCase{
		{
			name: "successful write with no headers",
			before: func(t *testing.T, tc *writeJSONTestCase) {
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			after: func(t *testing.T, tc *writeJSONTestCase) {
				// No cleanup needed for this test case
			},
			args: struct {
				status  int
				data    any
				headers []http.Header
			}{
				status: http.StatusOK,
				data:   map[string]string{"key": "value"},
			},
			expected: struct {
				status  int
				headers http.Header
				body    string
				err     bool
			}{
				status:  http.StatusOK,
				headers: http.Header{"Content-Type": []string{"application/json"}},
				body:    `{"key":"value"}` + "\n",
				err:     false,
			},
		},
		{
			name: "successful write with custom headers",
			before: func(t *testing.T, tc *writeJSONTestCase) {
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			after: func(t *testing.T, tc *writeJSONTestCase) {
				// No cleanup needed for this test case
			},
			args: struct {
				status  int
				data    any
				headers []http.Header
			}{
				status: http.StatusCreated,
				data:   map[string]string{"key": "value"},
				headers: []http.Header{{
					"X-Custom-Header": {"custom-value"},
				}},
			},
			expected: struct {
				status  int
				headers http.Header
				body    string
				err     bool
			}{
				status: http.StatusCreated,
				headers: http.Header{
					"Content-Type":    []string{"application/json"},
					"X-Custom-Header": []string{"custom-value"},
				},
				body: `{"key":"value"}` + "\n",
				err:  false,
			},
		},
		{
			name: "error on invalid JSON",
			before: func(t *testing.T, tc *writeJSONTestCase) {
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			after: func(t *testing.T, tc *writeJSONTestCase) {
				// No cleanup needed for this test case
			},
			args: struct {
				status  int
				data    any
				headers []http.Header
			}{
				status: http.StatusOK,
				data:   make(chan int), // Channels can't be JSON marshaled
			},
			expected: struct {
				status  int
				headers http.Header
				body    string
				err     bool
			}{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt
			tc.deps = &testDependencies{}

			// before
			if tc.before != nil {
				tc.before(t, &tc)
			}

			// Register after
			if tc.after != nil {
				t.Cleanup(func() {
					tc.after(t, &tc)
				})
			}

			// Execute
			err := WriteJSON(
				tc.deps.responseRecorder,
				tc.args.status,
				tc.args.data,
				tc.args.headers...,
			)

			// Verify
			if tc.expected.err {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Check status code
			assert.Equal(t, tc.expected.status, tc.deps.responseRecorder.Code)

			// Check headers
			for k, v := range tc.expected.headers {
				assert.Equal(t, v, tc.deps.responseRecorder.Header()[k])
			}

			// Check body
			if tc.expected.body != "" {
				var actual, expected any
				err = json.Unmarshal(tc.deps.responseRecorder.Body.Bytes(), &actual)
				require.NoError(t, err)
				err = json.Unmarshal([]byte(tc.expected.body), &expected)
				require.NoError(t, err)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

type writeJSONErrorTestCase struct {
	name   string
	deps   *testDependencies
	before func(t *testing.T, tc *writeJSONErrorTestCase)
	after  func(t *testing.T, tc *writeJSONErrorTestCase)
	args   struct {
		status  int
		message string
		headers []http.Header
	}
	expected struct {
		status  int
		headers http.Header
		resp    *JSONResponse
		err     bool
	}
}

func TestWriteJSONError(t *testing.T) {
	tests := []writeJSONErrorTestCase{
		{
			name: "successful error response",
			before: func(t *testing.T, tc *writeJSONErrorTestCase) {
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			after: func(t *testing.T, tc *writeJSONErrorTestCase) {
				// No cleanup needed for this test case
			},
			args: struct {
				status  int
				message string
				headers []http.Header
			}{
				status:  http.StatusBadRequest,
				message: "invalid input",
				headers: []http.Header{{"X-Request-ID": {"123"}}},
			},
			expected: struct {
				status  int
				headers http.Header
				resp    *JSONResponse
				err     bool
			}{
				status: http.StatusBadRequest,
				headers: http.Header{
					"Content-Type":  {"application/json"},
					"X-Request-ID": {"123"},
				},
				resp: &JSONResponse{
					Error:   true,
					Message: "invalid input",
				},
				err: false,
			},
		},
		{
			name: "empty error message",
			before: func(t *testing.T, tc *writeJSONErrorTestCase) {
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
				}
			},
			after: func(t *testing.T, tc *writeJSONErrorTestCase) {
				// No cleanup needed for this test case
			},
			args: struct {
				status  int
				message string
				headers []http.Header
			}{
				status:  http.StatusInternalServerError,
				message: "",
			},
			expected: struct {
				status  int
				headers http.Header
				resp    *JSONResponse
				err     bool
			}{
				status: http.StatusInternalServerError,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
				resp: &JSONResponse{
					Error:   true,
					Message: "",
				},
				err: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt
			tc.deps = &testDependencies{}

			// Setup
			if tc.before != nil {
				tc.before(t, &tc)
			}

			// Register after
			if tc.after != nil {
				t.Cleanup(func() {
					tc.after(t, &tc)
				})
			}

			// Execute
			err := WriteJSONError(
				tc.deps.responseRecorder,
				tc.args.status,
				tc.args.message,
				tc.args.headers...,
			)

			// Verify
			if tc.expected.err {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Check status code
			assert.Equal(t, tc.expected.status, tc.deps.responseRecorder.Code)

			// Check headers
			for k, v := range tc.expected.headers {
				assert.Equal(t, v, tc.deps.responseRecorder.Header()[k])
			}

			// Check response body
			var resp JSONResponse
			err = json.Unmarshal(tc.deps.responseRecorder.Body.Bytes(), &resp)
			require.NoError(t, err)
			assert.Equal(t, tc.expected.resp, &resp)
		})
	}
}

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type readJSONTestCase struct {
	name   string
	deps   *testDependencies
	before func(t *testing.T, tc *readJSONTestCase)
	after  func(t *testing.T, tc *readJSONTestCase)
	args   struct {
		target any
	}
	expected struct {
		result      any
		err         bool
		errContains string
	}
}

func TestReadJSON(t *testing.T) {
	tests := []readJSONTestCase{
		{
			name: "successful read",
			before: func(t *testing.T, tc *readJSONTestCase) {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"test","value":42}`))
				req.Header.Set("Content-Type", "application/json")
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			after: func(t *testing.T, tc *readJSONTestCase) {
				// No cleanup needed for this test case
			},
			args: struct{ target any }{
				target: &testStruct{},
			},
			expected: struct {
				result      any
				err         bool
				errContains string
			}{
				result: &testStruct{
					Name:  "test",
					Value: 42,
				},
				err: false,
			},
		},
		{
			name: "empty body",
			before: func(t *testing.T, tc *readJSONTestCase) {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))
				req.Header.Set("Content-Type", "application/json")
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			after: func(t *testing.T, tc *readJSONTestCase) {
				// No cleanup needed for this test case
			},
			args: struct{ target any }{
				target: &testStruct{},
			},
			expected: struct {
				result      any
				err         bool
				errContains string
			}{
				result:      nil,
				err:         true,
				errContains: "EOF",
			},
		},
		{
			name: "invalid JSON",
			before: func(t *testing.T, tc *readJSONTestCase) {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("{invalid-json}"))
				req.Header.Set("Content-Type", "application/json")
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			after: func(t *testing.T, tc *readJSONTestCase) {
				// No cleanup needed for this test case
			},
			args: struct{ target any }{
				target: &testStruct{},
			},
			expected: struct {
				result      any
				err         bool
				errContains string
			}{
				result:      nil,
				err:         true,
				errContains: "invalid character 'i' looking for beginning of object key string",
			},
		},
		{
			name: "unknown field",
			before: func(t *testing.T, tc *readJSONTestCase) {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"test","value":42,"unknown":"field"}`))
				req.Header.Set("Content-Type", "application/json")
				tc.deps = &testDependencies{
					responseRecorder: httptest.NewRecorder(),
					request:          req,
				}
			},
			after: func(t *testing.T, tc *readJSONTestCase) {
				// No cleanup needed for this test case
			},
			args: struct{ target any }{
				target: &testStruct{},
			},
			expected: struct {
				result      any
				err         bool
				errContains string
			}{
				result:      nil,
				err:         true,
				errContains: "json: unknown field",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt
			tc.deps = &testDependencies{}

			// Setup
			if tc.before != nil {
				tc.before(t, &tc)
			}

			// Register after
			if tc.after != nil {
				t.Cleanup(func() {
					tc.after(t, &tc)
				})
			}

			// Execute
			err := ReadJSON(tc.deps.responseRecorder, tc.deps.request, tc.args.target)

			// Verify
			if tc.expected.err {
				assert.Error(t, err)
				if tc.expected.errContains != "" {
					assert.Contains(t, err.Error(), tc.expected.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected.result, tc.args.target)

			// Check that the request had the correct content type
			assert.Equal(t, "application/json", tc.deps.request.Header.Get("Content-Type"))
		})
	}
}

func TestReadJSON_MaxBytes(t *testing.T) {
	type args struct {
		target any
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
