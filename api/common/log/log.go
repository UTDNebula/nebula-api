package log

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
)

var output = os.Stdout //zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

var Logger = zerolog.New(output).With().Timestamp().Logger()

func AddCodeLocation(e *zerolog.Event) *zerolog.Event {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	file := filepath.Base(frame.File)
	fn := filepath.Base(frame.Function)
	line := frame.Line
	return e.Str("file", file).Str("fn", fn).Int("line", line)
}

func WriteDebug(msg string) {
	Logger.Debug().Msg(msg)
}

func WriteError(err error) {
	AddCodeLocation(Logger.Err(err)).Send()
}

func WriteErrorMsg(msg string) {
	AddCodeLocation(Logger.Error()).Msg(msg)
}

func WriteErrorWithMsg(err error, msg string) {
	AddCodeLocation(Logger.Err(err)).Msg(msg)
}

func WritePanic(err error) {
	AddCodeLocation(Logger.Panic().Err(err)).Send()
}

var OffsetNotTypeInteger = "Offset is not type integer"

// Error should be used for Internal API errors. Any client side error should have a high-level log with Info
// or a low level log with Debug.

// Debug should be used for low-level logs, such as database transactions or anything detailing the API logic only.

// Info should be used for logging API requests at the high level.
