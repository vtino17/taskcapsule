package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "0.1.0"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func String() string {
	return fmt.Sprintf("taskcapsule %s\ncommit: %s\nbuilt: %s\ngo: %s", Version, Commit, BuildDate, runtime.Version())
}
