package log

import (
	"os"
	"runtime"

	"github.com/rs/zerolog"
)

var Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

func WriteError(err error) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()

	Logger.Err(err).Str("file", file).Str("fn", fn).Int("line", line).Send()
}

func WriteErrorWithMsg(err error, msg string) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()

	Logger.Err(err).Str("file", file).Str("fn", fn).Int("line", line).Msg(msg)
}

func WritePanic(err error) {
	Logger.Panic().Err(err).Send()
}

var OffsetNotTypeInteger = "Offset is not type integer"

// Error should be used for Internal API errors. Any client side error should have a high-level log with Info
// or a low level log with Debug.

// Debug should be used for low-level logs, such as database transactions or anything detailing the API logic only.

// Info should be used for logging API requests at the high level.
