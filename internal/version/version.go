// Package version handles the human read-able version for the consul-telemetry-collector itself and important
// components
package version

import (
	_ "embed" // import embed package to embed Version file
	"strings"
)

var (
	// GitCommit is the current Git SHA of the built project. It is set by LD flags.
	GitCommit string

	// Version is the static version tag of the project
	//go:embed VERSION
	Version string
)

// GetHumanVersion returns the human read-able version of the project version information.
func GetHumanVersion() string {
	version := Version

	return strings.TrimSpace(version)
}
