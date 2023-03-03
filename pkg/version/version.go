// Package version handles the human read-able version for the consul-telemetry-collector itself and important
// components
package version

import "fmt"

var (
	// GitCommit is the current Git SHA of the built project. It is set by LD flags
	GitCommit string

	// Version is the static version tag of the project
	Version = "0.0.1"

	// VersionPrerelease is the prerelease version string
	VersionPrerelease = "dev"
)

// GetHumanVersion returns the human read-able version of the project version information
func GetHumanVersion() string {
	version := Version
	release := VersionPrerelease

	return fmt.Sprintf("%s-%s", version, release)

}
