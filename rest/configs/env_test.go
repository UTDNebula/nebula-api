package configs

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

// The environment readers below exit through log.Fatalf when a required
// variable is missing, which would take the test binary down with them. Each
// case is instead run in a subprocess that calls the reader directly, so the
// exit status can be observed.
const fatalCaseVar = "NEBULA_ENV_FATAL_CASE"

var fatalCases = map[string]func(){
	"mongo_uri":   func() { GetEnvMongoURI() },
	"login_netid": func() { GetEnvLogin() },
	"login_password": func() {
		os.Setenv("LOGIN_NETID", "abc123456")
		GetEnvLogin()
	},
}

// unsetEnv removes a variable for the duration of the test, restoring whatever
// value it had once the test finishes.
func unsetEnv(t *testing.T, key string) {
	t.Helper()

	// t.Setenv records the original value so it is restored on cleanup; the
	// unset immediately afterwards is what the test actually wants.
	t.Setenv(key, "")
	os.Unsetenv(key)
}

// TestEnvFatalSubprocess is the entry point for the subprocesses spawned by
// runFatalCase. It does nothing during a normal test run.
func TestEnvFatalSubprocess(t *testing.T) {

	name := os.Getenv(fatalCaseVar)
	if name == "" {
		t.Skip("not running as a fatal-case subprocess")
	}

	fatalCase, ok := fatalCases[name]
	if !ok {
		t.Fatalf("unknown fatal case %q", name)
	}

	fatalCase()
}

// runFatalCase re-runs this test binary with only TestEnvFatalSubprocess
// enabled, and returns the resulting exit code.
func runFatalCase(t *testing.T, name string) int {
	t.Helper()

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("could not locate the test binary: %v", err)
	}

	cmd := exec.Command(executable, "-test.run", "^TestEnvFatalSubprocess$")
	cmd.Env = []string{fatalCaseVar + "=" + name}

	// Run from an empty directory so the subprocess cannot pick up a
	// contributor's .env while walking up from the package directory.
	cmd.Dir = t.TempDir()

	err = cmd.Run()

	exitError, isExitError := err.(*exec.ExitError)
	switch {
	case err == nil:
		return 0
	case isExitError:
		return exitError.ExitCode()
	default:
		t.Fatalf("could not run the fatal case %q: %v", name, err)
		return -1
	}
}

// captureLog redirects the standard logger for the duration of the test and
// returns what was written to it.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()

	var captured bytes.Buffer

	originalOutput := log.Writer()
	originalFlags := log.Flags()
	log.SetOutput(&captured)
	log.SetFlags(0)

	t.Cleanup(func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	})

	return &captured
}

// resetEnvWarnings clears the once-guards so a test can observe a warning that
// an earlier test has already consumed.
func resetEnvWarnings(t *testing.T) {
	t.Helper()

	limitWarnOnce = sync.Once{}
	uploadSizeWarnOnce = sync.Once{}
	uploadSizeCapWarnOnce = sync.Once{}
}

func TestGetPortString(t *testing.T) {

	t.Run("Unset returns the default port", func(t *testing.T) {

		unsetEnv(t, "PORT")

		if port := GetPortString(); port != ":8080" {
			t.Errorf("expected \":8080\" when PORT is unset, got %q", port)
		}
	})

	t.Run("Set returns the configured port", func(t *testing.T) {

		t.Setenv("PORT", "3000")

		if port := GetPortString(); port != ":3000" {
			t.Errorf("expected \":3000\", got %q", port)
		}
	})
}

func TestGetEnvMongoURI(t *testing.T) {

	t.Run("Set returns the URI", func(t *testing.T) {

		const expected = "mongodb://localhost:27017"
		t.Setenv("MONGODB_URI", expected)

		if uri := GetEnvMongoURI(); uri != expected {
			t.Errorf("expected %q, got %q", expected, uri)
		}
	})

	t.Run("Unset exits", func(t *testing.T) {

		if code := runFatalCase(t, "mongo_uri"); code != 1 {
			t.Errorf("expected exit code 1 when MONGODB_URI is unset, got %d", code)
		}
	})
}

func TestGetClubsDBUri(t *testing.T) {

	t.Run("Set returns the URI", func(t *testing.T) {

		const expected = "mongodb://localhost:27017/clubs"
		t.Setenv("CLUBS_DB_URI", expected)

		if uri := GetClubsDBUri(); uri != expected {
			t.Errorf("expected %q, got %q", expected, uri)
		}
	})

	t.Run("Unset panics", func(t *testing.T) {

		unsetEnv(t, "CLUBS_DB_URI")

		defer func() {
			if recovered := recover(); recovered == nil {
				t.Error("expected a panic when CLUBS_DB_URI is unset")
			}
		}()

		GetClubsDBUri()
	})
}

func TestGetEnvLogin(t *testing.T) {

	t.Run("Set returns both values", func(t *testing.T) {

		const (
			expectedNetID    = "abc123456"
			expectedPassword = "hunter2"
		)
		t.Setenv("LOGIN_NETID", expectedNetID)
		t.Setenv("LOGIN_PASSWORD", expectedPassword)

		netID, password := GetEnvLogin()

		if netID != expectedNetID {
			t.Errorf("expected net ID %q, got %q", expectedNetID, netID)
		}
		if password != expectedPassword {
			t.Errorf("expected password %q, got %q", expectedPassword, password)
		}
	})

	t.Run("Missing net ID exits", func(t *testing.T) {

		if code := runFatalCase(t, "login_netid"); code != 1 {
			t.Errorf("expected exit code 1 when LOGIN_NETID is unset, got %d", code)
		}
	})

	t.Run("Missing password exits", func(t *testing.T) {

		if code := runFatalCase(t, "login_password"); code != 1 {
			t.Errorf("expected exit code 1 when LOGIN_PASSWORD is unset, got %d", code)
		}
	})
}

func TestGetEnvLimit(t *testing.T) {

	const defaultLimit int64 = 20

	t.Run("Unset returns the default limit", func(t *testing.T) {

		unsetEnv(t, "LIMIT")

		if limit := GetEnvLimit(); limit != defaultLimit {
			t.Errorf("expected the default limit %d, got %d", defaultLimit, limit)
		}
	})

	t.Run("Set returns the parsed limit without warning", func(t *testing.T) {

		t.Setenv("LIMIT", "50")
		resetEnvWarnings(t)
		captured := captureLog(t)

		if limit := GetEnvLimit(); limit != 50 {
			t.Errorf("expected 50, got %d", limit)
		}

		if logged := captured.String(); logged != "" {
			t.Errorf("expected no log output for a valid LIMIT, got %q", logged)
		}
	})

	t.Run("Unparseable falls back to the default limit and warns", func(t *testing.T) {

		t.Setenv("LIMIT", "not-a-number")
		resetEnvWarnings(t)
		captured := captureLog(t)

		if limit := GetEnvLimit(); limit != defaultLimit {
			t.Errorf("expected the default limit %d, got %d", defaultLimit, limit)
		}

		logged := captured.String()
		for _, want := range []string{"LIMIT", "not-a-number", "20"} {
			if !strings.Contains(logged, want) {
				t.Errorf("expected the warning to mention %q, got %q", want, logged)
			}
		}
	})

	t.Run("Warns only once per process", func(t *testing.T) {

		t.Setenv("LIMIT", "not-a-number")
		resetEnvWarnings(t)
		captured := captureLog(t)

		GetEnvLimit()
		GetEnvLimit()
		GetEnvLimit()

		// GetEnvLimit runs on every paginated request, so a misconfigured
		// LIMIT must not log once per request.
		if lines := strings.Count(captured.String(), "\n"); lines != 1 {
			t.Errorf("expected exactly 1 warning across 3 calls, got %d: %q", lines, captured.String())
		}
	})
}

func TestGetEnvMaxUploadSize(t *testing.T) {

	const (
		defaultLimit int64 = 30 * 1024 * 1024
		hardCapLimit int64 = 50 * 1024 * 1024
	)

	t.Run("Unset returns the default size", func(t *testing.T) {

		unsetEnv(t, "MAX_UPLOAD_SIZE")

		if size := GetEnvMaxUploadSize(); size != defaultLimit {
			t.Errorf("expected the default size %d, got %d", defaultLimit, size)
		}
	})

	t.Run("Set returns the parsed size", func(t *testing.T) {

		t.Setenv("MAX_UPLOAD_SIZE", "100")

		if size := GetEnvMaxUploadSize(); size != 100 {
			t.Errorf("expected 100, got %d", size)
		}
	})

	t.Run("Unparseable falls back to the default size and warns", func(t *testing.T) {

		t.Setenv("MAX_UPLOAD_SIZE", "not-a-number")
		resetEnvWarnings(t)
		captured := captureLog(t)

		if size := GetEnvMaxUploadSize(); size != defaultLimit {
			t.Errorf("expected the default size %d, got %d", defaultLimit, size)
		}

		logged := captured.String()
		for _, want := range []string{"MAX_UPLOAD_SIZE", "not-a-number", "31457280"} {
			if !strings.Contains(logged, want) {
				t.Errorf("expected the warning to mention %q, got %q", want, logged)
			}
		}
	})

	t.Run("Above the hard cap returns the hard cap and warns", func(t *testing.T) {

		t.Setenv("MAX_UPLOAD_SIZE", "104857600")
		resetEnvWarnings(t)
		captured := captureLog(t)

		if size := GetEnvMaxUploadSize(); size != hardCapLimit {
			t.Errorf("expected the hard cap %d, got %d", hardCapLimit, size)
		}

		logged := captured.String()
		for _, want := range []string{"MAX_UPLOAD_SIZE", "104857600", "52428800"} {
			if !strings.Contains(logged, want) {
				t.Errorf("expected the warning to mention %q, got %q", want, logged)
			}
		}
	})
}

func TestGetEnvSentryDSN(t *testing.T) {

	t.Run("Unset returns empty string", func(t *testing.T) {

		unsetEnv(t, "SENTRY_DSN")

		if dsn := GetEnvSentryDSN(); dsn != "" {
			t.Errorf("expected an empty DSN when SENTRY_DSN is unset, got %q", dsn)
		}
	})

	t.Run("Set returns the DSN", func(t *testing.T) {

		const expected = "https://public@sentry.example.com/1"
		t.Setenv("SENTRY_DSN", expected)

		if dsn := GetEnvSentryDSN(); dsn != expected {
			t.Errorf("expected %q, got %q", expected, dsn)
		}
	})
}
