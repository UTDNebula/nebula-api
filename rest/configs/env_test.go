package configs

import (
	"os"
	"testing"
)

func TestGetEnvSentryDSN(t *testing.T) {

	original, hadOriginal := os.LookupEnv("SENTRY_DSN")
	defer func() {
		if hadOriginal {
			os.Setenv("SENTRY_DSN", original)
		} else {
			os.Unsetenv("SENTRY_DSN")
		}
	}()

	t.Run("Unset returns empty string", func(t *testing.T) {

		os.Unsetenv("SENTRY_DSN")

		if dsn := GetEnvSentryDSN(); dsn != "" {
			t.Errorf("expected an empty DSN when SENTRY_DSN is unset, got %q", dsn)
		}
	})

	t.Run("Set returns the DSN", func(t *testing.T) {

		const expected = "https://public@sentry.example.com/1"
		os.Setenv("SENTRY_DSN", expected)

		if dsn := GetEnvSentryDSN(); dsn != expected {
			t.Errorf("expected %q, got %q", expected, dsn)
		}
	})
}
