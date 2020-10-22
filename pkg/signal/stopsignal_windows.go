package signal

import (
	"os"
)

var shutdownSignals = []os.Signal{os.Interrupt}
