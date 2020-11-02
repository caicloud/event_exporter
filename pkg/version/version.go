package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "1.0.0"
	Branch    = "UNKNOWN"
	Commit    = "UNKNOWN"
	BuildUser = "Caicloud Authors"
	BuildDate = "UNKNOWN"
	GoVersion = runtime.Version()
)

func Message() string {
	const format = `event_exporter:%s (Branch: %s, Revision: %s)
build user: %s
build date: %s
go version: %s
version   : %s
`
	return fmt.Sprintf(format, Version, Branch, Commit, BuildUser, BuildDate, GoVersion, Version)
}
