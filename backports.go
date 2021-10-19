package v8go

import "time"

// Backport time.UnixMicro from go 1.17 - https://pkg.go.dev/time#UnixMicro
// timeUnixMicro accepts microseconds and converts to nanoseconds to be used
// with time.Unix which returns the local Time corresponding to the given Unix time,
// usec microseconds since January 1, 1970 UTC.
func timeUnixMicro(usec int64) time.Time {
	return time.Unix(0, usec*1000)
}
