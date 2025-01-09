package application_test

import (
	"fmt"
	"testing"

	"github.com/flohansen/auther/internal/application"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestSSLMode_UnmarshalYAML(t *testing.T) {
	t.Run("should default to require", func(t *testing.T) {
		var sslmode application.SSLMode
		assert.Equal(t, application.SSLModeRequire, sslmode)
	})

	t.Run("should unmarshal", func(t *testing.T) {
		tests := []struct {
			input    string
			expected application.SSLMode
		}{
			{input: "disable", expected: application.SSLModeDisable},
			{input: "require", expected: application.SSLModeRequire},
			{input: "verify-ca", expected: application.SSLModeVerifyCA},
			{input: "verify-full", expected: application.SSLModeVerifyFull},
		}

		for _, test := range tests {
			// given
			var sslmode application.SSLMode

			// when
			err := yaml.Unmarshal([]byte(test.input), &sslmode)

			// then
			assert.NoError(t, err)
			assert.Equal(t, test.expected, sslmode)
		}
	})
}

func TestSSLMode_String(t *testing.T) {
	// given
	tests := []struct {
		input    application.SSLMode
		expected string
	}{
		{input: application.SSLModeDisable, expected: "disable"},
		{input: application.SSLModeRequire, expected: "require"},
		{input: application.SSLModeVerifyCA, expected: "verify-ca"},
		{input: application.SSLModeVerifyFull, expected: "verify-full"},
	}

	for _, test := range tests {
		// when
		out := fmt.Sprintf("%s", test.input)

		// then
		assert.Equal(t, test.expected, out)
	}
}
