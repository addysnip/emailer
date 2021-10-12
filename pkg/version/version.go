package version

import "fmt"

var (
	Version   = "0.99"
	GitCommit = "HEAD"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func FriendlyVersion() string {
	return fmt.Sprintf("%s-%s (built: %s [GoVersion %s])", Version, GitCommit, BuildTime, GoVersion)
}
