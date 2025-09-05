package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getStringTestCase struct {
	name   string
	before func(t *testing.T, tc *getStringTestCase)
	after  func(t *testing.T, tc *getStringTestCase)
	args   struct {
		key      string
		fallback string
	}
	expected struct {
		result string
		err    error
	}
}

func TestGetString(t *testing.T) {
	tests := []getStringTestCase{
		{
			name: "should return environment variable when it exists",
			before: func(t *testing.T, tc *getStringTestCase) {
				t.Setenv(tc.args.key, "test_value")
			},
			after: func(t *testing.T, tc *getStringTestCase) {
				t.Cleanup(func() {
					os.Unsetenv(tc.args.key)
				})
			},
			args: struct {
				key      string
				fallback string
			}{
				key:      "TEST_KEY",
				fallback: "default_value",
			},
			expected: struct {
				result string
				err    error
			}{
				result: "test_value",
			},
		},
		{
			name: "should return fallback when environment variable does not exist",
			before: func(t *testing.T, tc *getStringTestCase) {
				// No setup needed for this test case
			},
			after: func(t *testing.T, tc *getStringTestCase) {
				// No teardown needed for this test case
			},
			args: struct {
				key      string
				fallback string
			}{
				key:      "NON_EXISTENT_KEY",
				fallback: "default_value",
			},
			expected: struct {
				result string
				err    error
			}{
				result: "default_value",
			},
		},
		{
			name: "should handle empty environment variable",
			before: func(t *testing.T, tc *getStringTestCase) {
				t.Setenv(tc.args.key, "")
			},
			after: func(t *testing.T, tc *getStringTestCase) {
				t.Cleanup(func() {
					os.Unsetenv(tc.args.key)
				})
			},
			args: struct {
				key      string
				fallback string
			}{
				key:      "EMPTY_KEY",
				fallback: "default_value",
			},
			expected: struct {
				result string
				err    error
			}{
				result: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test case
			tc := tt

			// Setup
			if tc.before != nil {
				tc.before(t, &tc)
			}

			// Register teardown
			if tc.after != nil {
				t.Cleanup(func() {
					tc.after(t, &tc)
				})
			}

			// Execute
			result := GetString(tc.args.key, tc.args.fallback)

			// Verify
			assert.Equal(t, tc.expected.result, result)
			if tc.expected.err != nil {
				assert.ErrorIs(t, tc.expected.err, tc.expected.err)
			} else {
				assert.NoError(t, tc.expected.err)
			}
		})
	}
}

type getIntTestCase struct {
	name   string
	before func(t *testing.T, tc *getIntTestCase)
	after  func(t *testing.T, tc *getIntTestCase)
	args   struct {
		key      string
		fallback int
	}
	expected struct {
		result int
		err    error
	}
}

func TestGetInt(t *testing.T) {
	tests := []getIntTestCase{
		{
			name: "should return parsed integer when environment variable is a valid integer",
			before: func(t *testing.T, tc *getIntTestCase) {
				t.Setenv(tc.args.key, "42")
			},
			after: func(t *testing.T, tc *getIntTestCase) {
				t.Cleanup(func() {
					os.Unsetenv(tc.args.key)
				})
			},
			args: struct {
				key      string
				fallback int
			}{
				key:      "TEST_INT",
				fallback: 0,
			},
			expected: struct {
				result int
				err    error
			}{
				result: 42,
			},
		},
		{
			name: "should return fallback when environment variable is not a valid integer",
			before: func(t *testing.T, tc *getIntTestCase) {
				t.Setenv(tc.args.key, "not_an_int")
			},
			after: func(t *testing.T, tc *getIntTestCase) {
				t.Cleanup(func() {
					os.Unsetenv(tc.args.key)
				})
			},
			args: struct {
				key      string
				fallback int
			}{
				key:      "INVALID_INT",
				fallback: 100,
			},
			expected: struct {
				result int
				err    error
			}{
				result: 100,
			},
		},
		{
			name: "should return fallback when environment variable does not exist",
			before: func(t *testing.T, tc *getIntTestCase) {
				// No before needed
			},
			after: func(t *testing.T, tc *getIntTestCase) {
				// No after needed
			},
			args: struct {
				key      string
				fallback int
			}{
				key:      "NON_EXISTENT_INT",
				fallback: 200,
			},
			expected: struct {
				result int
				err    error
			}{
				result: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			if tc.before != nil {
				tc.before(t, &tc)
			}

			if tc.after != nil {
				t.Cleanup(func() {
					tc.after(t, &tc)
				})
			}

			result := GetInt(tc.args.key, tc.args.fallback)

			assert.Equal(t, tc.expected.result, result)
			if tc.expected.err != nil {
				assert.ErrorIs(t, tc.expected.err, tc.expected.err)
			} else {
				assert.NoError(t, tc.expected.err)
			}
		})
	}
}

type getBoolTestCase struct {
	name   string
	before func(t *testing.T, tc *getBoolTestCase)
	after  func(t *testing.T, tc *getBoolTestCase)
	args   struct {
		key      string
		fallback bool
	}
	expected struct {
		result bool
		err    error
	}
}

func TestGetBool(t *testing.T) {
	tests := []getBoolTestCase{
		{
			name: "should return true when environment variable is 'true' (case insensitive)",
			before: func(t *testing.T, tc *getBoolTestCase) {
				t.Setenv(tc.args.key, "true")
			},
			after: func(t *testing.T, tc *getBoolTestCase) {
				t.Cleanup(func() {
					os.Unsetenv(tc.args.key)
				})
			},
			args: struct {
				key      string
				fallback bool
			}{
				key:      "BOOL_TRUE",
				fallback: false,
			},
			expected: struct {
				result bool
				err    error
			}{
				result: true,
			},
		},
		{
			name: "should return false when environment variable is 'false' (case insensitive)",
			before: func(t *testing.T, tc *getBoolTestCase) {
				t.Setenv(tc.args.key, "FALSE")
			},
			after: func(t *testing.T, tc *getBoolTestCase) {
				t.Cleanup(func() {
					os.Unsetenv(tc.args.key)
				})
			},
			args: struct {
				key      string
				fallback bool
			}{
				key:      "BOOL_FALSE",
				fallback: true,
			},
			expected: struct {
				result bool
				err    error
			}{
				result: false,
			},
		},
		{
			name: "should return fallback when environment variable is not a valid boolean",
			before: func(t *testing.T, tc *getBoolTestCase) {
				t.Setenv(tc.args.key, "not_a_boolean")
			},
			after: func(t *testing.T, tc *getBoolTestCase) {
				t.Cleanup(func() {
					os.Unsetenv(tc.args.key)
				})
			},
			args: struct {
				key      string
				fallback bool
			}{
				key:      "INVALID_BOOL",
				fallback: true,
			},
			expected: struct {
				result bool
				err    error
			}{
				result: true,
			},
		},
		{
			name: "should return fallback when environment variable does not exist",
			before: func(t *testing.T, tc *getBoolTestCase) {
				// No before needed
			},
			after: func(t *testing.T, tc *getBoolTestCase) {
				// No after needed
			},
			args: struct {
				key      string
				fallback bool
			}{
				key:      "NON_EXISTENT_BOOL",
				fallback: false,
			},
			expected: struct {
				result bool
				err    error
			}{
				result: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			if tc.before != nil {
				tc.before(t, &tc)
			}

			if tc.after != nil {
				t.Cleanup(func() {
					tc.after(t, &tc)
				})
			}

			result := GetBool(tc.args.key, tc.args.fallback)

			assert.Equal(t, tc.expected.result, result)
			if tc.expected.err != nil {
				assert.ErrorIs(t, tc.expected.err, tc.expected.err)
			} else {
				assert.NoError(t, tc.expected.err)
			}
		})
	}
}
