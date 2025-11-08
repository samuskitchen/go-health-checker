package zerolog

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	rslog "github.com/rs/zerolog/log"
)

// captureStderr redirects os.Stderr to a pipe, executes action(), and returns
// what was written to stderr.
func captureStderr(action func()) (string, error) {
	// Guardar y restaurar la salida estándar de error
	originalStderr := os.Stderr
	defer func() { os.Stderr = originalStderr }()

	// Create the pipe to capture stderr
	reader, writer, pipeErr := os.Pipe()
	if pipeErr != nil {
		return "", pipeErr
	}

	// Redirect stderr to the writer
	os.Stderr = writer

	// Execute the action whose output errors we want to capture
	action()

	// Close the writer and check for errors
	if closeErr := writer.Close(); closeErr != nil {
		return "", closeErr
	}

	// Read from the reader and check for errors
	var buffer bytes.Buffer
	if _, copyErr := io.Copy(&buffer, reader); copyErr != nil {
		return "", copyErr
	}

	// Return what was captured
	return buffer.String(), nil
}

func TestInitLogger(t *testing.T) {
	t.Run("GlobalLevel_without_debug", func(t *testing.T) {
		// We force a different level at the beginning
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		InitLogger("appX", false)
		got := zerolog.GlobalLevel()
		if got != zerolog.InfoLevel {
			t.Errorf("Expected global level = InfoLevel; got %v", got)
		}
	})

	t.Run("GlobalLevel_with_debug", func(t *testing.T) {
		// We restart
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		InitLogger("appY", true)
		got := zerolog.GlobalLevel()
		if got != zerolog.DebugLevel {
			t.Errorf("Expected global level = DebugLevel; got %v", got)
		}
	})

	t.Run("Output_includes_debug_message_and_app_field", func(t *testing.T) {
		out, err := captureStderr(func() {
			InitLogger("myApp", true)
			// InitLogger now outputs "Debug mode enabled" if debug=true
			// Add another debug as test mode
			rslog.Debug().Msg("unit-test")
		})
		if err != nil {
			t.Fatalf("Error capturing stderr: %v", err)
		}

		if !strings.Contains(out, "Debug mode enabled") {
			t.Error("Se esperaba ver 'Debug mode enabled' en la salida")
		}

		if !strings.Contains(out, "| unit-test |") {
			t.Error("Se esperaba ver el mensaje '| unit-test |' formateado")
		}

		if !strings.Contains(out, "app:myApp") {
			t.Error("Se esperaba el campo 'app=myApp' en la salida")
		}
	})
}
