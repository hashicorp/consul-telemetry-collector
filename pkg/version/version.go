package version

import "fmt"

var (
	GitCommit string

	Version = "0.0.1"

	VersionPrerelease = "dev"
)

func GetHumanVersion() string {
	version := Version
	release := VersionPrerelease

	return fmt.Sprintf("%s-%s", version, release)
}
