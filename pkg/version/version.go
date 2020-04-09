package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "UNKNOWN"
	Branch    = "UNKNOWN"
	Commit    = "UNKNOWN"
	BuildUser = "Caicloud Authors"
	BuildDate = "UNKNOWN"
)

func Message() string {
	const format = `event_exporter:%s (Branch: %s, Revision: %s)
  build user: %s
  build date: %s
  go version: %s
`
	return fmt.Sprintf(format, Version, Branch, Commit, BuildUser, BuildDate, runtime.Version())
}
