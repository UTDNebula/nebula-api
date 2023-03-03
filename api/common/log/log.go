package log

import (
	"log"
	"os"
)

var (
	// Error should be used for Internal API errors. Any client side error should have a high-level log with Info
	// or a low level log with Debug.
	Error = log.New(os.Stderr, "[Error] ", log.Lmsgprefix)

	// Debug should be used for low-level logs, such as database transactions or anything detailing the API logic only.
	Debug = log.New(os.Stderr, "[Debug] ", log.Lmsgprefix)

	// Info should be used for logging API requests at the high level.
	Info = log.New(os.Stdout, "[Info]  ", log.Lmsgprefix)
)
